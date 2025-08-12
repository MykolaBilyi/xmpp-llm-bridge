package middleware

import (
	"__template__/internal/ports"
	"__template__/pkg/middleware"
	"context"
	"errors"
	"net/http"
	"time"
)

type TimeoutMiddleware struct {
	timeout        time.Duration
	loggerProvider ports.LoggerProvider
	next           http.Handler
}

func (m TimeoutMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger := m.loggerProvider.Value(req.Context())
	ctx, cancel := context.WithTimeout(req.Context(), m.timeout)
	defer cancel() // Always do this to clean up contexts, otherwise they'll hang out and gather since they are blocked go rountines

	m.next.ServeHTTP(w, req.WithContext(ctx))

	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Warn("request timeout")
			http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		}
	}
}

func NewTimeoutMiddleware(config ports.Config, loggerProvider ports.LoggerProvider) middleware.Middleware {
	config.SetDefault("timeout", "30s")

	return func(next http.Handler) http.Handler {
		return &TimeoutMiddleware{
			timeout:        config.GetDuration("timeout"),
			loggerProvider: loggerProvider,
			next:           next,
		}
	}
}
