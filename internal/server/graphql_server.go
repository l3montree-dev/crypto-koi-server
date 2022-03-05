package server

import (
	"context"
	"fmt"
	"image"
	"strconv"

	imageDraw "golang.org/x/image/draw"

	"image/draw"
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
	generator         generator.Generator
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
					http_util.WriteHttpError(w, http.StatusInternalServerError, "invalid token: %e", err)
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					// Token is either expired or not active yet
					s.logger.Infof("token is either expired or not active yet: %s", err)
					http_util.WriteHttpError(w, http.StatusUnauthorized, "token is either expired or not active yet: %e", err)
				} else {
					s.logger.Errorf("invalid token: %s", err)
					http_util.WriteHttpError(w, http.StatusInternalServerError, "invalid token: %e", err)
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
	preloader := generator.NewMemoryPreloader(imagesBasePath)
	return &GraphqlServer{db: db, generator: generator.NewGenerator(preloader), logger: orchardclient.Logger.WithField("component", "GraphqlServer")}
}

func (s *GraphqlServer) imageHandlerFactory(size int, drawBackgroundColor bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenId := chi.URLParam(r, "tokenId")
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

		img, koi := s.generator.TokenId2Image(tokenId)

		scaledImg := image.NewRGBA(image.Rect(0, 0, size, size))

		if drawBackgroundColor {
			primaryColor := koi.GetAttributes().PrimaryColor
			for y := 0; y < scaledImg.Bounds().Max.Y; y++ {
				for x := 0; x < scaledImg.Bounds().Max.X; x++ {
					scaledImg.Set(x, y, primaryColor)
				}
			}
		}

		imageDraw.BiLinear.Scale(scaledImg, scaledImg.Rect, img, img.Bounds(), draw.Over, nil)

		w.Header().Set("Content-Type", "image/png")

		png.Encode(w, scaledImg)
	}
}

func (s *GraphqlServer) getLeaderboardUpdateRoutine() leader.Listener {
	sleepTime := os.Getenv("LEADERBOARD_UPDATE_INTERVAL")
	if sleepTime == "" {
		sleepTime = fmt.Sprint(60 * 5)
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
				s.logger.WithField("took", time.Since(now).String()).Infof("leaderboard update routine finished. Err: %v", err)
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
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"*"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

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
	userSvc := service.NewUserService(userRepository)
	authSvc := service.NewAuthService(userRepository, tokenSvc)
	eventSvc := service.NewEventService(eventRepository)
	gameSvc := service.NewGameService(gameRepository, eventSvc, tokenSvc)
	// init all controllers
	cryptogotchiSvc := service.NewCryptogotchiService(cryptogotchiRepository)
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
	router.Route("/auth", func(r chi.Router) {
		// make sure to stop processing after 10 seconds.
		r.Use(middleware.Timeout(10 * time.Second))
		r.Post("/login", authController.Login)
		r.Post("/refresh", authController.Refresh)
	})

	router.Route("/v1", func(r chi.Router) {
		// make sure to stop processing after 10 seconds.
		r.Use(middleware.Timeout(10 * time.Second))
		// opensea.io integration.
		// gets called by their API and wallet applications.
		r.Get("/tokens/{tokenId}", openseaController.GetCryptogotchi)
		r.Get("/images/{tokenId}", s.imageHandlerFactory(350, false))
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
	resolver := graph.NewResolver(int(chainId), s.userSvc, eventSvc, cryptogotchiSvc, gameSvc, authSvc, cryptokoiApi, s.generator)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}))

	// authorized routes
	router.Group(func(r chi.Router) {
		// make sure to stop processing after 10 seconds.
		r.Use(graphqlTimeout(10 * time.Second))
		// attach the auth middleware to the router
		r.Use(s.authMiddleware)
		r.Handle("/query", srv)
	})

	s.logger.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
