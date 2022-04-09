package server

import (
	"context"
	"fmt"
	"image"
	"strconv"

	"image/png"
	"log"
	"net/http"

	"os"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/graph"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/graph/generated"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/config"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/controller"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/cryptokoi"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/generator"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/http_util"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/repositories"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/service"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/pkg/leader"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
	"gorm.io/gorm"
)

type GraphqlServer struct {
	db                *gorm.DB
	tokenSvc          service.TokenSvc
	userSvc           service.UserSvc
	cryptogotchiSvc   service.CryptogotchiSvc
	koiGenerator      generator.Generator
	dragonGenerator   generator.Generator
	leaderElection    leader.LeaderElection
	cryptokoiListener *cryptokoi.CryptoKoiEventListener
	logger            *logrus.Entry
}

type responseData struct {
	status int
	size   int
}
type loggingResponseWriter struct {
	http.ResponseWriter // compose original http.ResponseWriter
	responseData        *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
	r.responseData.size += size            // capture size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
	r.responseData.status = statusCode       // capture status code
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		responseData := &responseData{}
		loggingWriter := loggingResponseWriter{ResponseWriter: w, responseData: responseData}
		next.ServeHTTP(&loggingWriter, r)
		orchardclient.Logger.
			WithField("method", r.Method).
			WithField("status", loggingWriter.responseData.status).
			WithField("size", loggingWriter.responseData.size).
			WithField("path", r.URL.Path).
			WithField("took", time.Since(now).String()).Info("handled request")
	})
}

// the auth middleware will set the current logged in user into the context
func (s *GraphqlServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the token from the header
		token := r.Header.Get("Authorization")
		if token == "" {
			s.logger.Warn("auth middleware called without token")
			next.ServeHTTP(w, r)
			// http_util.WriteHttpError(w, http.StatusUnauthorized, "no token provided")
			return
		}

		// parse the token from the request
		// remove the bearer prefix
		claims, err := s.tokenSvc.ParseToken(strings.Replace(strings.Replace(token, "Bearer ", "", -1), "bearer ", "", -1))
		if err != nil {
			// invalid token
			// log it.
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					s.logger.Errorf("invalid token: %s", err)
					http_util.WriteHttpError(w, http.StatusInternalServerError, fmt.Sprintf("invalid token: %e", err))
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					// Token is either expired or not active yet
					s.logger.Infof("token is either expired or not active yet: %s", err)
					http_util.WriteHttpError(w, http.StatusUnauthorized, fmt.Sprintf("token is either expired or not active yet: %e", err))
				} else {
					s.logger.Errorf("invalid token: %s", err)
					http_util.WriteHttpError(w, http.StatusInternalServerError, fmt.Sprintf("invalid token: %e", err))
				}
			}
			return
		}

		// get the user from the token
		user, err := s.userSvc.GetById(claims.(jwt.MapClaims)["sub"].(string))
		if err != nil {
			// invalid user
			// log it.
			s.logger.Errorf("invalid user: %s", err)
			http_util.WriteHttpError(w, http.StatusInternalServerError, "could not fetch user from database")
			return
		}

		oldCtx := r.Context()
		newCtx := context.WithValue(oldCtx, config.USER_CTX_KEY, &user)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func graphqlTimeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer func() {
				cancel()
			}()

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
func NewGraphqlServer(db *gorm.DB, imagesBasePath string) Server {
	koiPreloader := generator.NewMemoryPreloader(imagesBasePath + "/koi")
	dragonPreloader := generator.NewMemoryPreloader(imagesBasePath + "/dragon")

	// build the caches during bootstrap in a non blocking way
	go koiPreloader.BuildCachesForSizes([]int{200, 350, 1024})
	go dragonPreloader.BuildCachesForSizes([]int{200, 350, 1024})

	return &GraphqlServer{
		db:              db,
		koiGenerator:    generator.NewGenerator(koiPreloader),
		dragonGenerator: generator.NewGenerator(dragonPreloader),
		logger:          orchardclient.Logger.WithField("component", "GraphqlServer"),
	}
}

func (s *GraphqlServer) imageHandlerFactory(defaultSize int, drawBackgroundColor bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenId := chi.URLParam(r, "tokenId")
		t := r.URL.Query().Get("type")
		if t == "" || !strings.EqualFold("dragon", t) {
			t = "koi"
		}
		// check if hex.
		if strings.IndexFunc(tokenId, util.IsNotDigit) > -1 {
			// not only digits - use as hex.
			tmp, err := util.UuidToUint256(tokenId)
			if err != nil {
				http_util.WriteHttpError(w, http.StatusBadRequest, "invalid tokenId")
				return
			}

			tokenId = tmp.String()
		}

		var img image.Image
		var koi *cryptokoi.CryptoKoi

		var err error
		size := defaultSize
		if sz := r.URL.Query().Get("size"); sz != "" {
			size, err = strconv.Atoi(sz)
			if err != nil {
				http_util.WriteHttpError(w, http.StatusBadRequest, "invalid size")
				return
			}

		}

		// decide wether to generate a koi or a dragon
		if t == "koi" {
			img, koi = s.koiGenerator.TokenId2Image(tokenId, size)
		} else {
			img, koi = s.dragonGenerator.TokenId2Image(tokenId, size)
		}

		if drawBackgroundColor {
			primaryColor := koi.GetAttributes().PrimaryColor
			for y := 0; y < img.Bounds().Max.Y; y++ {
				for x := 0; x < img.Bounds().Max.X; x++ {
					img.(*image.RGBA).Set(x, y, primaryColor)
				}
			}
		}

		w.Header().Set("Content-Type", "image/png")

		png.Encode(w, img)
	}
}

func (s *GraphqlServer) getLeaderboardUpdateRoutine() leader.Listener {
	sleepTime := os.Getenv("LEADERBOARD_UPDATE_INTERVAL")
	if sleepTime == "" {
		sleepTime = fmt.Sprint(30)
	}
	sleepTimeInt, err := strconv.Atoi(sleepTime)
	orchardclient.FailOnError(err, "could not parse leaderboard update interval")
	return leader.NewListener(func(cancelChan <-chan struct{}) {
		for {
			select {
			case <-cancelChan:
				return
			default:
				now := time.Now()
				s.logger.Info("leaderboard update routine started")
				err := s.cryptogotchiSvc.UpdateRanks()
				if err != nil {
					s.logger.WithField("took", time.Since(now).String()).Errorf("leaderboard update routine finished. Err: %v", err)
				} else {
					s.logger.WithField("took", time.Since(now).String()).Info("leaderboard update routine finished")
				}

				time.Sleep(time.Second * time.Duration(sleepTimeInt))
			}
		}
	})
}

func (s *GraphqlServer) getBlockchainListener() leader.Listener {
	return leader.NewListener(func(cancelChan <-chan struct{}) {
		eventChan := s.cryptokoiListener.StartListener()
		for {
			select {
			case <-cancelChan:
				return
			case ev := <-eventChan:
				crypt, err := s.cryptogotchiSvc.GetCryptogotchiByUint256(ev.TokenId)

				if err != nil {
					s.logger.Error(err)
					continue
				}

				err = s.cryptogotchiSvc.MarkAsNft(&crypt)
				if err != nil {
					s.logger.Error(err)
					continue
				}
			}
		}
	})
}

func (s *GraphqlServer) getLeaderElection() leader.LeaderElection {
	// create new leader election object to make sure, that we run the listener only once - even in a distributed environment.
	podName := os.Getenv("POD_NAME")
	var leaderElection leader.LeaderElection
	if podName != "" {
		namespace := os.Getenv("NAMESPACE")
		if namespace == "" {
			namespace = "default"
		}
		// distributed environment detected.
		leaderElection = leader.NewKubernetesLeaderElection(context.Background(), "leaderelection", namespace)
	} else {
		// local environment detected.
		leaderElection = leader.NewAlwaysLeader()
	}
	return leaderElection
}

func (s *GraphqlServer) Start() {
	defaultPort := "8080"
	isDev := os.Getenv("DEV") != ""
	port := os.Getenv("PORT")

	imageBaseUrl := os.Getenv("IMAGE_BASE_URL")
	if imageBaseUrl == "" {
		s.logger.Fatal("IMAGE_BASE_URL env variable is not set.")
	}

	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})
	// register all middlewares
	router.Use(loggerMiddleware)

	// allow cross origin request
	router.Use(cors.AllowAll().Handler)

	router.Use(middleware.Recoverer)
	router.Use(middleware.AllowContentType("application/json"))

	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)

	router.Use(sentryMiddleware.Handle)

	router.Get("/images/{tokenId}", s.imageHandlerFactory(1024, false))
	router.Get("/thumbnails/{tokenId}", s.imageHandlerFactory(200, false))

	// init all repositories
	cryptogotchiRepository := repositories.NewGormCryptogotchiRepository(s.db)
	eventRepository := repositories.NewGormEventRepository(s.db)
	userRepository := repositories.NewGormUserRepository(s.db)
	gameRepository := repositories.NewGormGameStatRepository(s.db)

	// init all services
	tokenSvc := service.NewTokenService()
	apiKey := os.Getenv("FCM_API_KEY")
	if apiKey == "" {
		orchardclient.Logger.Panic("FCM_API_KEY is not defined")
	}

	notificationSvc := service.NewNotificationSvc(apiKey)
	userSvc := service.NewUserService(userRepository)
	authSvc := service.NewAuthService(userRepository, tokenSvc)
	eventSvc := service.NewEventService(eventRepository)
	gameSvc := service.NewGameService(gameRepository, eventSvc, tokenSvc)
	// init all controllers
	cryptogotchiSvc := service.NewCryptogotchiService(cryptogotchiRepository, userRepository, notificationSvc)
	authController := controller.NewAuthController(userRepository, cryptogotchiSvc, authSvc)
	openseaController := controller.NewOpenseaController(imageBaseUrl, eventRepository, cryptogotchiSvc)

	// set services to server instance for middleware and listeners
	s.tokenSvc = tokenSvc
	s.userSvc = userSvc
	s.cryptogotchiSvc = cryptogotchiSvc

	// add all routes.
	if isDev {
		router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}
	router.Group(func(r chi.Router) {
		// make sure to stop processing after 10 seconds.
		r.Use(middleware.Timeout(10 * time.Second))
		r.Post("/auth/login", authController.Login)
		r.Post("/auth/register", authController.Register)
		r.Post("/auth/refresh", authController.Refresh)
	})

	router.Route("/v1", func(r chi.Router) {
		// make sure to stop processing after 10 seconds.
		r.Use(middleware.Timeout(10 * time.Second))
		// opensea.io integration.
		// gets called by their API and wallet applications.
		r.Get("/tokens/{tokenId}", openseaController.GetCryptogotchi)
		r.Get("/images/{tokenId}", s.imageHandlerFactory(350, false))
		r.Get("/fakes/{tokenId}", openseaController.GetFakeCryptogotchi)
	})

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		s.logger.Fatal("PRIVATE_KEY environment variable is not defined")
	}

	chainUrl := os.Getenv("CHAIN_URL")
	if chainUrl == "" {
		s.logger.Fatal("CHAIN_URL environment variable is not defined")
	}
	chainWs := os.Getenv("CHAIN_WS")
	if chainWs == "" {
		s.logger.Fatal("CHAIN_WS environment variable is not defined")
	}

	ethHttpClient, err := ethclient.Dial(chainUrl)
	if err != nil {
		s.logger.Fatal(err)
	}
	defer ethHttpClient.Close()

	ethWsClient, err := ethclient.Dial(chainWs)
	if err != nil {
		s.logger.Fatal(err)
	}
	defer ethWsClient.Close()

	contractAddress := os.Getenv("CONTRACT_ADDRESS")
	if contractAddress == "" {
		s.logger.Fatal("CONTRACT_ADDRESS is not set")
	}

	httpBinding, err := cryptokoi.NewCryptoKoiBinding(common.HexToAddress(contractAddress), ethHttpClient)
	orchardclient.FailOnError(err, "Failed to instantiate a CryptoKoi contract binding (HTTP)")
	wsBinding, err := cryptokoi.NewCryptoKoiBinding(common.HexToAddress(contractAddress), ethWsClient)
	orchardclient.FailOnError(err, "Failed to instantiate a CryptoKoi contract binding (WS)")

	cryptokoiApi := cryptokoi.NewCryptokoiApi(privateKey, httpBinding)
	s.cryptokoiListener = cryptokoi.NewCryptoKoiEventListener(wsBinding)

	s.leaderElection = s.getLeaderElection()
	// start the listener.
	s.leaderElection.AddListener(s.getBlockchainListener())
	s.leaderElection.AddListener(s.getLeaderboardUpdateRoutine())
	s.leaderElection.AddListener(cryptogotchiSvc.GetNotificationListener())
	// start all listeners
	go s.leaderElection.RunElection()

	chainIdEnv := os.Getenv("CHAIN_ID")
	if chainIdEnv == "" {
		s.logger.Fatal("CHAIN_ID is not set")
	}
	chainId, err := strconv.ParseInt(chainIdEnv, 10, 64)
	if err != nil {
		s.logger.Fatal(err)
	}

	// attach the graphql handler to the router
	resolver := graph.NewResolver(int(chainId), s.userSvc, eventSvc, cryptogotchiSvc, gameSvc, authSvc, cryptokoiApi, s.koiGenerator)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}))

	// authorized routes
	router.Group(func(r chi.Router) {
		// make sure to stop processing after 10 seconds.
		r.Use(graphqlTimeout(10 * time.Second))
		// attach the auth middleware to the router
		r.Use(s.authMiddleware)
		r.Handle("/query", srv)
	})

	router.Group(func(r chi.Router) {
		r.Use(s.authMiddleware)
		r.Delete("/auth", authController.DestroyAccount)
	})

	s.logger.Infof("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
