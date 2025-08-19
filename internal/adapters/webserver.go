package adapters

import (
	"context"
	"net"
	"net/http"

	"xmpp-llm-bridge/internal/ports"
	"xmpp-llm-bridge/internal/providers"
)

type WebServer struct {
	config         ports.Config
	loggerProvider *providers.LoggerProvider
	server         *http.Server
}

var _ ports.Server = (*WebServer)(nil)

func NewWebServer(ctx context.Context, config ports.Config, loggerProvider *providers.LoggerProvider, handler http.Handler) (*WebServer, error) {
	config.SetDefault("listenAddr", ":8080")

	return &WebServer{
		config:         config,
		loggerProvider: loggerProvider,
		server: &http.Server{
			Addr:        config.GetString("listenAddr"),
			Handler:     handler,
			BaseContext: func(_ net.Listener) context.Context { return ctx },
		},
	}, nil
}

func (s *WebServer) Serve(ctx context.Context) error {
	logger := s.loggerProvider.Value(ctx)
	l, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	logger.Info("web server started", ports.Fields{"address": s.server.Addr})

	return s.server.Serve(l)
}

func (s *WebServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
