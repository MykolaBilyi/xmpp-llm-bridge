package main

import (
	"context"
	"errors"
	"net/http"
	"syscall"

	"xmpp-llm-bridge/api/client"
	"xmpp-llm-bridge/api/server"
	"xmpp-llm-bridge/internal/adapters"
	"xmpp-llm-bridge/internal/ports"

	"log"
	"os"
	"os/signal"
)

func main() {
	config, err := adapters.NewConfig()
	if err != nil {
		log.Printf("error reading config: %v\n", err)
		os.Exit(1)
	}

	logger, err := adapters.NewLogger(config.Sub("logger"))
	if err != nil {
		log.Printf("error creating logger: %v\n", err)
		os.Exit(1)
	}

	handler := client.NewHandler(config.Sub("llm"), logger)
	jabber, err := adapters.NewJabberClient(config.Sub("xmpp"), handler, logger)

	if err != nil {
		logger.Error("error registering client", ports.Fields{"error": err})
		os.Exit(1)
	}

	go func() {
		if err := jabber.Serve(); err != nil {
			logger.Error("xmpp client error", ports.Fields{"error": err})
			os.Exit(1)
		}
	}()

	router := server.NewRouter(config.Sub("middleware"), logger)
	server := adapters.NewWebServer(config.Sub("http"), router, logger)

	go func() {
		if err := server.Serve(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			logger.Error("web server error", ports.Fields{"error": err})
			os.Exit(1)
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	<-signalChan
	logger.Info("shutting down")

	go func() {
		<-signalChan
		// second signal. Exit directly.
		logger.Warn("terminating")
		os.Exit(1)
	}()

	// graceful shutdown
	gracefulCtx, cancelShutdown := context.WithTimeout(context.Background(), config.GetDuration("application.shutdownTimeout"))
	defer cancelShutdown() // release resources afterwards

	if err := jabber.Shutdown(gracefulCtx); err != nil {
		logger.Error("shutdown error", ports.Fields{"error": err})
		defer os.Exit(1)
		return
	}

	if err := server.Shutdown(gracefulCtx); err != nil {
		logger.Error("shutdown error", ports.Fields{"error": err})
		defer os.Exit(1)
		return
	}

	defer os.Exit(0)
}
