package main

import (
	"context"
	"syscall"

	"xmpp-llm-bridge/internal/adapters"
	"xmpp-llm-bridge/internal/app"

	"log"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()

	appConfig, err := adapters.NewConfig(ctx)
	if err != nil {
		log.Printf("error reading config: %v\n", err)
		os.Exit(1)
	}

	application, err := app.NewApplication(
		ctx,
		appConfig,
	)
	if err != nil {
		log.Printf("error initialization application: %v\n", err)
		os.Exit(1)
	}

	err = application.Start(ctx)
	if err != nil {
		log.Printf("error starting application: %v\n", err)
		os.Exit(1)
	}

	stopCh, closeCh := createChannel()
	defer closeCh()

	<-stopCh
	go terminateOnSecondSignal(stopCh)

	err = application.Stop(ctx)
	if err != nil {
		log.Printf("error during shutdown: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func createChannel() (chan os.Signal, func()) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return stopCh, func() {
		close(stopCh)
	}
}

func terminateOnSecondSignal(stopCh chan os.Signal) {
	<-stopCh
	os.Exit(1)
}
