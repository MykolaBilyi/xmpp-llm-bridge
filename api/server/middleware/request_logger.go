package middleware

import (
	"__template__/internal/ports"
	"__template__/pkg/middleware"
	"context"
	"net/http"
	"time"

	"logur.dev/logur"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	if lrw.wroteHeader {
		return
	}
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
	lrw.wroteHeader = true
}

type RequestLoggerMiddleware struct {
	loggerProvider ports.LoggerProvider
	logLevel       logur.Level
	next           http.Handler
}

func (m *RequestLoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := m.loggerProvider.Value(r.Context())
	log := logur.LevelFunc(logger, m.logLevel)
	log("request", ports.Fields{"method": r.Method, "path": r.URL.Path})
	lrw := &loggingResponseWriter{ResponseWriter: w}

	defer func(start time.Time) {
		log("response", ports.Fields{"status": lrw.statusCode, "duration": time.Since(start).String()})
	}(time.Now())
	m.next.ServeHTTP(lrw, r)
}

func NewRequestLoggerMiddleware(config ports.Config, loggerProvider ports.LoggerProvider) middleware.Middleware {
	config.SetDefault("level", "debug")
	logLevel, ok := logur.ParseLevel(config.GetString("level"))

	if !ok {
		loggerProvider.Value(context.Background()).
			Warn("invalid log level for request logging, using info", ports.Fields{"value": config.GetString("level")})
		logLevel = logur.Info
	}

	return func(next http.Handler) http.Handler {
		return &RequestLoggerMiddleware{
			loggerProvider: loggerProvider,
			logLevel:       logLevel,
			next:           next,
		}
	}
}
