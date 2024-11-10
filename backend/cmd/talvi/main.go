package main

import (
	"context"
	"log"
	"net/http"

	"github.com/alpineworks/ootel"
	"github.com/gorilla/mux"
	"github.com/michaelpeterswa/talvi/backend/internal/accounts"
	"github.com/michaelpeterswa/talvi/backend/internal/cockroach"
	"github.com/michaelpeterswa/talvi/backend/internal/config"
	"github.com/michaelpeterswa/talvi/backend/internal/dragonfly"
	"github.com/michaelpeterswa/talvi/backend/internal/encryption"
	"github.com/michaelpeterswa/talvi/backend/internal/handlers"
	"github.com/michaelpeterswa/talvi/backend/internal/logging"
	"github.com/michaelpeterswa/talvi/backend/internal/middleware"
	"github.com/michaelpeterswa/talvi/backend/internal/movies"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	k, err := config.Get("talvi")
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}

	logger, err := logging.InitLogging(k.String(config.LogLevel))
	if err != nil {
		log.Fatalf("error initializing logging: %v", err)
	}

	logger.Info("welcome to talvi backend!")

	ootelClient := ootel.NewOotelClient(
		ootel.WithMetricConfig(ootel.NewMetricConfig(
			k.Bool(config.MetricsEnabled),
			k.Int(config.MetricsPort),
		)),
		ootel.WithTraceConfig(ootel.NewTraceConfig(
			k.Bool(config.TracingEnabled),
			k.Float64(config.TracingRatio),
			k.String(config.ServiceName),
			k.String(config.ServiceVersion),
		)),
	)

	shutdown, err := ootelClient.Init(ctx)
	if err != nil {
		logger.Fatal("error initializing ootel client", zap.Error(err))
	}
	defer func() {
		_ = shutdown(ctx)
	}()

	logger.Info("ootel initialized")

	aesConfig, err := encryption.NewAESConfig(k.String(config.AESKey))
	if err != nil {
		logger.Fatal("error initializing aes config", zap.Error(err))
	}

	aesClient, err := encryption.NewAESClient(aesConfig)
	if err != nil {
		logger.Fatal("error initializing aes client", zap.Error(err))
	}

	dragonflyClient, err := dragonfly.NewDragonflyClient(k.String(config.DragonflyHost), k.Int(config.DragonflyPort), k.String(config.DragonflyAuth))
	if err != nil {
		logger.Fatal("error initializing dragonfly client", zap.Error(err))
	}

	if err := redisotel.InstrumentMetrics(dragonflyClient.GetClient(), redisotel.WithMeterProvider(otel.GetMeterProvider())); err != nil {
		logger.Fatal("failed to instrument redis metrics", zap.Error(err))
	}

	if err := redisotel.InstrumentTracing(dragonflyClient.GetClient()); err != nil {
		logger.Fatal("failed to instrument redis tracing", zap.Error(err))
	}

	cockroachClient, err := cockroach.NewCockroachClient(k.String(config.CockroachURL))
	if err != nil {
		logger.Fatal("error initializing cockroach client", zap.Error(err))
	}
	defer cockroachClient.Client.Close()

	jwtMiddlewareClient, err := middleware.NewJWTMiddleware(logger, k.String(config.JWESecret))
	if err != nil {
		logger.Fatal("error initializing jwt middleware client", zap.Error(err))
	}

	accountAuthorizationMiddleware := middleware.NewAccountAuthorizationMiddleware(logger)

	accountsClient := accounts.NewAccountsClient(dragonflyClient, cockroachClient, aesClient)
	accountsHandler := handlers.NewAccountsHandler(logger, accountsClient)

	moviesClient := movies.NewMoviesClient(dragonflyClient, cockroachClient)
	moviesHandler := handlers.NewMovieHandler(moviesClient)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()
	v1Router := apiRouter.PathPrefix("/v1").Subrouter()

	accountsRouter := v1Router.PathPrefix("/accounts").Subrouter()
	accountsRouter.HandleFunc("/account", accountsHandler.CreateAccount).Methods("POST")
	accountsRouter.HandleFunc("/account", accountsHandler.GetAccount).Methods("GET")

	twofactorRouter := accountsRouter.PathPrefix("/2fa").Subrouter()
	twofactorRouter.HandleFunc("/2fa", accountsHandler.Create2FA).Methods("POST")
	twofactorRouter.HandleFunc("/2fa", accountsHandler.Get2FA).Methods("GET")
	twofactorRouter.HandleFunc("/2fa", accountsHandler.Delete2FA).Methods("DELETE")
	twofactorRouter.HandleFunc("/verify", accountsHandler.Verify2FA).Methods("GET")
	twofactorRouter.HandleFunc("/generate", accountsHandler.Generate2FA).Methods("GET")
	twofactorRouter.HandleFunc("/validate", accountsHandler.Validate2FA).Methods("GET")

	v1Router.HandleFunc("/movies", moviesHandler.GetMovies)

	router.Use(
		// requests tracer middleware
		otelmux.Middleware(k.String(config.ServiceName)),
		// requests counter middleware
		middleware.RequestsCounterMiddleware(),
		// jwt middleware
		jwtMiddlewareClient.JWTMiddleware,
		// account authorization middleware
		accountAuthorizationMiddleware.IsAccountAuthorized,
		// cors middleware
		middleware.CORS,
	)

	corsWrappedRouter := cors.AllowAll().Handler(router)
	err = http.ListenAndServe(":8080", corsWrappedRouter)
	if err != nil {
		panic(err)
	}
}
