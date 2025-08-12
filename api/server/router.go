package server

import (
	"__template__/api/server/handlers"
	"__template__/api/server/middleware"
	"__template__/internal/ports"
	"__template__/internal/providers"
	"net/http"

	"github.com/go-chi/httprate"
	"github.com/gorilla/mux"
)

func NewRouter(config ports.Config, logger ports.Logger) http.Handler {
	router := mux.NewRouter()

	loggerProvider := providers.NewLoggerProvider(logger)
	requestIdProvider := providers.NewRequestIdProvider()

	router.Use(
		middleware.NewRequestIdMiddleware(config.Sub("requestID"), requestIdProvider),
		middleware.NewLoggerMiddleware(config.Sub("logger"), loggerProvider, requestIdProvider),
		middleware.NewSkipperMiddleware(
			config.Sub("requestLogger"),
			middleware.NewRequestLoggerMiddleware(config.Sub("requestLogger"), loggerProvider),
		),
		middleware.NewRecoverMiddleware(config.Sub("recover"), loggerProvider),
		middleware.NewSkipperMiddleware(
			config.Sub("rateLimiter"),
			httprate.LimitByIP(
				config.GetInt("rateLimiter.requests"),
				config.GetDuration("rateLimiter.duration"),
			),
		),
		middleware.NewTimeoutMiddleware(config.Sub("timeout"), loggerProvider),
		middleware.NewSkipperMiddleware(
			config.Sub("cors"),
			middleware.NewCORSMiddleware(config.Sub("cors")),
		),
	)
	router.Handle("/health", handlers.NewHealthCheck()).Methods("GET")

	return router
}
