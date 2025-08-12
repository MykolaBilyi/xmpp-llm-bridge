package main

import (
	"__template__/api/server"
	"__template__/internal/adapters"
	"__template__/internal/ports"
	"context"
	"syscall"

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

	router := server.NewRouter(config.Sub("middleware"), logger)
	server := adapters.NewWebServer(config.Sub("http"), router, logger)

	go func() {
		if err := server.Serve(); err != nil {
			logger.Error("server error", ports.Fields{"error": err})
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

	if err := server.Shutdown(gracefulCtx); err != nil {
		logger.Error("shutdown error", ports.Fields{"error": err})
		defer os.Exit(1)
		return
	}

	defer os.Exit(0)
}
