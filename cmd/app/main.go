package main

import (
	"context"
	"deez-nats/config"
	"deez-nats/internal/models"
	"deez-nats/pkg/logging"
	"deez-nats/pkg/storage/postgres"
	"fmt"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var log = logging.GetLogger()

	ctx := context.Background()

	var cfg config.Config
	err := confita.NewLoader(env.NewBackend()).Load(ctx, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.ConnectDB(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	err = db.AutoMigrate(&models.Order{})
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
