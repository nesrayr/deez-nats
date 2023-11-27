package publisher

import (
	"deez-nats/pkg/logging"
	"github.com/nats-io/stan.go"
)

type IPublisher interface {
	PublishData(data []byte, subject string) error
	Close() error
}

type Publisher struct {
	sc stan.Conn
	l  *logging.Logger
}

// Connect to NATS Streaming server
func NewPublisher(cfg Config, l logging.Logger) (*Publisher, error) {
	sc, err := stan.Connect(cfg.ClusterID, cfg.PublisherClient)
	if err != nil {
		return nil, err
	}

	return &Publisher{sc: sc, l: &l}, nil
}

func (p *Publisher) PublishData(data []byte, subject string) error {
	err := p.sc.Publish(subject, data)
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
