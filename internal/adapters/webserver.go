package adapters

import (
	"context"
	"net"
	"net/http"

	"xmpp-llm-bridge/internal/ports"
)

type WebServer struct {
	baseCtx       context.Context
	cancelContext context.CancelFunc
	cfg           ports.Config
	logger        ports.Logger
	server        *http.Server
}

var _ ports.Server = (*WebServer)(nil)

func NewWebServer(cfg ports.Config, router http.Handler, logger ports.Logger) *WebServer {
	cfg.SetDefault("listenAddr", ":8080")
	baseCtx, cancel := context.WithCancel(context.Background())

	return &WebServer{
		cfg:           cfg,
		baseCtx:       baseCtx,
		cancelContext: cancel,
		logger:        logger,
		server: &http.Server{
			Addr:        cfg.GetString("listenAddr"),
			Handler:     router,
			BaseContext: func(_ net.Listener) context.Context { return baseCtx },
		},
	}
}

func (s *WebServer) Serve() error {
	l, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	s.logger.Info("web server started", ports.Fields{"address": s.server.Addr})

	return s.server.Serve(l)
}

func (s *WebServer) Shutdown(ctx context.Context) error {
	s.cancelContext()
	return s.server.Shutdown(ctx)
}
