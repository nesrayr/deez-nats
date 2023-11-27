package migrations

import (
	"deez-nats/pkg/logging"
	"fmt"
	"github.com/pressly/goose"
	"os"
)

func MigrateDB(command string, log logging.Logger, dir string, arguments ...string) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
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
