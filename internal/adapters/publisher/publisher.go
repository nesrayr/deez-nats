package publisher

import (
	"deez-nats/pkg/logging"
	"encoding/json"
	"github.com/nats-io/stan.go"
)

type IPublisher interface {
	PublishData(data map[string]interface{}, subject string) error
	Close() error
}

type Publisher struct {
	sc stan.Conn
	l  *logging.Logger
}

// Connect to NATS Streaming server
func NewPublisher(cfg Config, l logging.Logger) (*Publisher, error) {
	l.Debug(cfg.ClusterID, cfg.PublisherClient)
	sc, err := stan.Connect(cfg.ClusterID, cfg.PublisherClient, stan.NatsURL(cfg.URL))
	if err != nil {
		return nil, err
	}

	return &Publisher{sc: sc, l: &l}, nil
}

func (p *Publisher) PublishData(data map[string]interface{}, subject string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		p.l.Errorf("error marshaling JSON: %v", err)
		return err
	}

	err = p.sc.Publish(subject, jsonData)
	if err != nil {
		p.l.Errorf("error publishing message: %v", err)
		return err
	}
	p.l.Info("data successfully published")

	return nil
}

func (p *Publisher) Close() error {
	return p.sc.Close()
}
