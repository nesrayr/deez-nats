package config

import (
	"deez-nats/internal/service/publisher"
	"deez-nats/internal/service/subscriber"
	"deez-nats/pkg/storage/postgres"
)

type Config struct {
	Host       string `config:"HOST" yaml:"host"`
	Port       string `config:"PORT" yaml:"port"`
	Database   postgres.Config
	Publisher  publisher.Config
	Subscriber subscriber.Config
}
