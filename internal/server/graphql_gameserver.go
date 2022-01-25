package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v4"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/graph"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/graph/generated"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/config"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/controller"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/http_util"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
	resolver "gitlab.com/l3montree/cryptogotchi/clodhopper/internal/resolvers"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/service"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
	"gorm.io/gorm"
)

type GraphqlGameserver struct {
	db       *gorm.DB
	tokenSvc service.TokenSvc
	userSvc  service.UserSvc
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
func (s *GraphqlGameserver) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the token from the header
		token := r.Header.Get("Authorization")
		if token == "" {
			orchardclient.Logger.Warn("auth middleware called without token")
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
					orchardclient.Logger.Errorf("invalid token: %s", err)
					http_util.WriteHttpError(w, http.StatusInternalServerError, "invalid token: %e", err)
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					// Token is either expired or not active yet
					orchardclient.Logger.Infof("token is either expired or not active yet: %s", err)
					http_util.WriteHttpError(w, http.StatusUnauthorized, "token is either expired or not active yet: %e", err)
				} else {
					orchardclient.Logger.Errorf("invalid token: %s", err)
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
			orchardclient.Logger.Errorf("invalid user: %s", err)
			next.ServeHTTP(w, r)
			http_util.WriteHttpError(w, http.StatusInternalServerError, "could not fetch user from database")
			return
		}

		oldCtx := r.Context()
		newCtx := context.WithValue(oldCtx, config.USER_CTX_KEY, &user)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func NewGraphqlGameserver(db *gorm.DB) Server {
	return &GraphqlGameserver{db: db}
}

func (s *GraphqlGameserver) Start() {
	defaultPort := "8080"
	isDev := os.Getenv("DEV") != ""
	port := os.Getenv("PORT")

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
	// make sure to stop processing after 30 seconds.
	router.Use(middleware.Timeout(30 * time.Second))
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)

	router.Use(sentryMiddleware.Handle)

	tokenSvc := service.NewTokenService()
	s.tokenSvc = tokenSvc

	cryptogotchiRepository := repositories.NewGormCryptogotchiRepository(s.db)
	eventRepository := repositories.NewGormEventRepository(s.db)
	userRepository := repositories.NewGormUserRepository(s.db)

	authController := controller.NewAuthController(userRepository, tokenSvc)

	openseaController := controller.NewOpenseaController(eventRepository, cryptogotchiRepository)

	s.userSvc = service.NewUserService(userRepository)

	if isDev {
		router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", authController.Login)
		r.Post("/refresh", authController.Refresh)
	})

	router.Route("/integrations", func(r chi.Router) {
		// opensea.io integration.
		// gets called by their API and wallet applications.
		r.Get("/opensea/:tokenId", openseaController.GetCryptogotchi)
	})

	cryptogotchiResolver := resolver.NewCryptogotchiResolver(eventRepository, cryptogotchiRepository)

	// attach the graphql handler to the router
	resolver := graph.NewResolver(&cryptogotchiResolver)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}))

	// authorized routes
	router.Group(func(r chi.Router) {
		// attach the auth middleware to the router
		r.Use(s.authMiddleware)
		r.Handle("/query", srv)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
