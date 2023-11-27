package main

import (
	"context"
	"deez-nats/config"
	"deez-nats/internal/adapters/subscriber"
	"deez-nats/internal/migrations"
	"deez-nats/internal/repo"
	"deez-nats/pkg/logging"
	"deez-nats/pkg/storage/cache"
	"deez-nats/pkg/storage/postgres"
	"fmt"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"sync"
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

	log.Debug(cfg)

	db, err := postgres.ConnectDB(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrations.MigrateDB("up", log, "./internal/migrations")

	c := cache.NewCache()

	repository := repo.NewRepository(db, c, log)

	sub, err := subscriber.NewSubscriber(cfg.Subscriber, log, repository)
	if err != nil {
		log.Error(err)
	}
	defer func() {
		_ = sub.Close()
	}()
	log.Debug("successfully connected to nats-streaming")

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		sub.ReadMessages("subject")
	}()

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

	wg.Wait()
}
