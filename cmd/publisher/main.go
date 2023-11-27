package main

import (
	"context"
	"deez-nats/config"
	"deez-nats/internal/adapters/publisher"
	router "deez-nats/internal/ports/publisher"
	"deez-nats/pkg/logging"
	"fmt"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var log = logging.GetLogger()

	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Error(err)
	}

	var cfg config.Config
	err = confita.NewLoader(env.NewBackend()).Load(ctx, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	pub, err := publisher.NewPublisher(cfg.Publisher, log)
	if err != nil {
		log.Error(err)
	}
	defer func() {
		_ = pub.Close()
	}()

	r := router.SetupRoutes(pub, log)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r)
	if err != nil {
		log.Fatal(err)
	}

	// configuring graceful shutdown
	sigQuit := make(chan os.Signal, 1)
	defer close(sigQuit)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})
}
