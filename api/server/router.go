package server

import (
	"net/http"

	"xmpp-llm-bridge/api/server/handlers"
	"xmpp-llm-bridge/internal/ports"
)

func NewRouter(config ports.Config, logger ports.Logger) http.Handler {
	router := http.NewServeMux()

	router.Handle("GET /health", handlers.NewHealthCheck())

	return router
}
