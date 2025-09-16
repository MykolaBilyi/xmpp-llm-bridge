package handlers

import (
	"net/http"
	"xmpp-llm-bridge/internal/providers"
)

func NewHealthCheck(loggerProvider *providers.LoggerProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := loggerProvider.Value(r.Context())
		logger.Debug("health check")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK\n"))
	})
}
