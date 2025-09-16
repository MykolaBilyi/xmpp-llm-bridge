package app

import (
	"context"
	"errors"
	"net/http"
	"xmpp-llm-bridge/internal/adapters"
	"xmpp-llm-bridge/internal/app/client"
	"xmpp-llm-bridge/internal/app/server"
	"xmpp-llm-bridge/internal/ports"
	"xmpp-llm-bridge/internal/providers"
)

type Application struct {
	config ports.Config
	server ports.Server
	client ports.LLMService
	jabber ports.XMPPSession
	logger ports.Logger
}

func NewApplication(ctx context.Context, config ports.Config) (*Application, error) {
	logger, err := adapters.NewLogger(config.Sub("log"))
	if err != nil {
		return nil, err
	}
	loggerProvider := providers.NewLoggerProvider(logger)
	ctx = loggerProvider.WithLogger(ctx, logger)

	llmService, err := adapters.NewOpenAIClient(ctx, config.Sub("openai"), loggerProvider)
	if err != nil {
		return nil, err
	}

	xmppSession, err := adapters.NewJabberClient(ctx, config.Sub("xmpp"), loggerProvider)
	if err != nil {
		return nil, err
	}
	if err := xmppSession.Connect(ctx); err != nil {
		return nil, err
	}

	httpHandler := server.NewRouter(config.Sub("http"), loggerProvider)
	livecheckServer, err := adapters.NewWebServer(
		ctx,
		config.Sub("http"),
		loggerProvider,
		httpHandler,
	)
	if err != nil {
		return nil, err
	}

	return &Application{
		config: config,
		server: livecheckServer,
		jabber: xmppSession,
		client: llmService,
		logger: logger,
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	errChan := make(chan error, 2)

	loggerProvider := providers.NewLoggerProvider(a.logger)
	requestIdProvider := providers.NewRequestIdProvider()
	xmppHandler := client.NewHandler(
		a.config,
		loggerProvider,
		requestIdProvider,
		a.jabber,
		a.client,
	)

	go func() {
		if err := a.jabber.Handle(ctx, xmppHandler); err != nil {
			a.logger.Error("xmpp client error", ports.Fields{"error": err})
			errChan <- err
		}
	}()

	go func() {
		if err := a.server.Serve(ctx); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			a.logger.Error("web server error", ports.Fields{"error": err})
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (a *Application) Stop(ctx context.Context) error {
	gracefulCtx, cancelShutdown := context.WithTimeout(
		ctx,
		a.config.GetDuration("application.shutdownTimeout"),
	)
	defer cancelShutdown() // release resources afterwards

	errChan := make(chan error, 2)

	go func() {
		if err := a.jabber.Close(gracefulCtx); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	go func() {
		if err := a.server.Shutdown(gracefulCtx); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	var errors []error
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}
