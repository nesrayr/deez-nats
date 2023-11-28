package migrations

import (
	"deez-nats/pkg/logging"
	"deez-nats/pkg/storage/postgres"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func MigrateDB(command string, log logging.Logger, cfg postgres.Config, dir string, arguments ...string) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		"localhost", cfg.User, cfg.Password, cfg.Database, cfg.Port,
	)
	db, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		log.Fatalf("goose: failed to open DB:%v\n", err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalf("goose : failed to close connection: %v\n", err)
		}
	}()

	if err = goose.Run(command, db, dir, arguments...); err != nil {
		log.Fatalf("goose run failed: %v\n", err)
	}
}
