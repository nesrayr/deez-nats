package config

import (
	"deez-nats/internal/service/publisher"
	"deez-nats/internal/service/subscriber"
	"deez-nats/pkg/storage/postgres"
)

type Config struct {
	Host           string `config:"HOST" yaml:"host"`
	SubscriberPort string `config:"SUBSCRIBER_PORT" yaml:"subscriber_port"`
	PublisherPort  string `config:"PUBLISHER_PORT" yaml:"publisher_port"`
	Database       postgres.Config
	Publisher      publisher.Config
	Subscriber     subscriber.Config
}
