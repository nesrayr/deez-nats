package main

import (
	"context"
	"deez-nats/config"
	"deez-nats/internal/migrations"
	router "deez-nats/internal/ports/subscriber"

	"deez-nats/internal/repo"
	"deez-nats/internal/service/subscriber"
	"deez-nats/pkg/logging"
	"deez-nats/pkg/storage/cache"
	"deez-nats/pkg/storage/postgres"
	"fmt"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"net/http"
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

	db, err := postgres.ConnectDB(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrations.MigrateDB("up", log, cfg.Database, "./internal/migrations")

	c := cache.NewCache()

	repository := repo.NewRepository(db, c, log)

	err = repository.UploadCache(ctx)
	if err != nil {
		log.Fatal(err)
	}

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
		sub.ReadMessages(ctx, cfg.Publisher.Subject, &wg)
	}()

	r := router.SetupRoutes(repository, log)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.SubscriberPort), r)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

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
