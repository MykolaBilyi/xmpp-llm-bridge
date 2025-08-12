package middleware

import (
	"__template__/internal/ports"
	"__template__/internal/providers"
	"__template__/pkg/middleware"
	"net/http"
)

type LoggerMiddleware struct {
	key               string
	loggerProvider    ports.LoggerProvider
	requestIdProvider providers.RequestIdProvider
	next              http.Handler
}

func (m *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := m.loggerProvider.Value(ctx).
		WithFields(ports.Fields{
			m.key: *m.requestIdProvider.Value(ctx),
		})
	m.next.ServeHTTP(w, r.WithContext(m.loggerProvider.WithLogger(ctx, logger)))
}

func NewLoggerMiddleware(config ports.Config, loggerProvider ports.LoggerProvider, requestIdProvider providers.RequestIdProvider) middleware.Middleware {
	config.SetDefault("requestIdKey", "request_id")

	return func(next http.Handler) http.Handler {
		return &LoggerMiddleware{
			key:               config.GetString("requestIdKey"),
			loggerProvider:    loggerProvider,
			requestIdProvider: requestIdProvider,
			next:              next,
		}
	}
}
