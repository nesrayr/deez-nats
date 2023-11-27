package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type Storage struct {
	DB *gorm.DB
}

func ConnectDB(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	log.Println("Connected!")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running migrations..")
	//err = db.AutoMigrate(&models.Person{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
