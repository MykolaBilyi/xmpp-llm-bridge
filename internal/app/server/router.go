package server

import (
	"net/http"

	"xmpp-llm-bridge/internal/app/server/handlers"
	"xmpp-llm-bridge/internal/ports"
	"xmpp-llm-bridge/internal/providers"
)

func NewRouter(config ports.Config, loggerProvider *providers.LoggerProvider) http.Handler {
	router := http.NewServeMux()

	router.Handle("GET /health", handlers.NewHealthCheck(loggerProvider))

	return router
}
