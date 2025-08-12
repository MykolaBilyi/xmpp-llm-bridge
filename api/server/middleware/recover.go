package middleware

import (
	"__template__/internal/ports"
	"__template__/pkg/middleware"
	"net/http"
	"runtime/debug"
)

type RecoverMiddleware struct {
	stackTrace     bool
	loggerProvider ports.LoggerProvider
	next           http.Handler
}

func (m *RecoverMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := m.loggerProvider.Value(r.Context())
	defer func() {
		if err := recover(); err != nil {
			details := ports.Fields{"error": err}
			if m.stackTrace {
				details["stack"] = string(debug.Stack())
			}
			logger.Error("panic recovered", details)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()
	m.next.ServeHTTP(w, r)
}

func NewRecoverMiddleware(config ports.Config, loggerProvider ports.LoggerProvider) middleware.Middleware {
	config.SetDefault("stackTrace", false)

	return func(next http.Handler) http.Handler {
		return &RecoverMiddleware{
			stackTrace:     config.GetBool("stackTrace"),
			loggerProvider: loggerProvider,
			next:           next,
		}
	}
}
