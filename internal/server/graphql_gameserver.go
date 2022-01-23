package server

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v4"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/graph"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/graph/generated"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/config"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/controller"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/http_util"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/service"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
	"gorm.io/gorm"
)

type GraphqlGameserver struct {
	db       *gorm.DB
	tokenSvc service.TokenSvc
	userSvc  service.UserSvc
}

// the auth middleware will set the current logged in user into the context
func (s *GraphqlGameserver) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the token from the header
		token := r.Header.Get("Authorization")
		if token == "" {
			orchardclient.Logger.Errorf("auth middleware called without token")
			http_util.WriteHttpError(w, http.StatusUnauthorized, "no token provided")
			return
		}

		// parse the token from the request
		claims, err := s.tokenSvc.ParseToken(token)
		if err != nil {
			// invalid token
			// log it.
			orchardclient.Logger.Errorf("invalid token: %s", err)
			http_util.WriteHttpError(w, http.StatusUnauthorized, "invalid token: %e")
			return
		}

		// get the user from the token
		user, err := s.userSvc.GetById(claims.(jwt.RegisteredClaims).Subject)
		if err != nil {
			// invalid user
			// log it.
			orchardclient.Logger.Errorf("invalid user: %s", err)
			http_util.WriteHttpError(w, http.StatusInternalServerError, "could not fetch user from database")
			return
		}

		oldCtx := r.Context()
		newCtx := context.WithValue(oldCtx, config.USER_CTX_KEY, user)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func NewGraphqlGameserver(db *gorm.DB) Server {
	return &GraphqlGameserver{db: db}
}

func (s *GraphqlGameserver) Start() {
	defaultPort := "8080"

	router := chi.NewRouter()
	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(sentryMiddleware.Handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	cryptogotchiRepository := repositories.NewGormCryptogotchiRepository(s.db)
	eventRepository := repositories.NewGormEventRepository(s.db)
	userRepository := repositories.NewGormUserRepository(s.db)

	authController := controller.NewAuthController(userRepository)
	// cryptogotchiController := controller.NewCryptogotchiController(eventRepository, cryptogotchiRepository)
	openseaController := controller.NewOpenseaController(eventRepository, cryptogotchiRepository)

	// register all middlewares
	// allow cross origin request
	// TODO: ADD cors middleware

	// register the controller
	router.Post("/auth/login", authController.Login)
	router.Post("/auth/refresh", authController.Refresh)

	// opensea.io integration.
	// gets called by their API and wallet applications.
	router.Get("/integrations/opensea/:tokenId", openseaController.GetCryptogotchi)

	// attach the graphql handler to the router
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	// attach the auth middleware to the router
	router.Use(s.authMiddleware)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
