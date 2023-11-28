package subscriber

import (
	"context"
	"deez-nats/internal/models"
	"deez-nats/internal/repo"
	"deez-nats/pkg/logging"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"sync"
)

type ISubscriber interface {
	ReadMessages(subject string) error
	Close() error
}

type Subscriber struct {
	sc   stan.Conn
	repo repo.IRepository
	l    *logging.Logger
}

// Connect to NATS Streaming server
func NewSubscriber(cfg Config, l logging.Logger, repo repo.IRepository) (*Subscriber, error) {
	sc, err := stan.Connect(cfg.ClusterID, cfg.SubscriberClient, stan.NatsURL(cfg.URL))
	if err != nil {
		return nil, err
	}

	return &Subscriber{sc: sc, l: &l, repo: repo}, nil
}

func (s *Subscriber) ReadMessages(ctx context.Context, subject string, wg *sync.WaitGroup) {
	s.l.Debug("reading messages")
	wg.Add(1)
	sub, _ := s.sc.Subscribe(subject, func(msg *stan.Msg) {
		defer wg.Done()
		var receivedOrder models.Order
		err := json.Unmarshal(msg.Data, &receivedOrder)
		if err != nil {
			s.l.Errorf("cannot unmarshal data %v", err)
		}
		s.l.Debug(receivedOrder)
		err = s.repo.AddOrder(ctx, receivedOrder)
		if err != nil {
			s.l.Errorf("error adding order %v", err)
		}
	}, stan.DeliverAllAvailable())
	wg.Wait()
	defer func(sub stan.Subscription) {
		s.l.Debug("closing subscriber")
		err := sub.Close()
		if err != nil {
			s.l.Errorf("error closing subscriber %v", err)
		}
	}(sub)
}

func (s *Subscriber) Close() error {
	return s.sc.Close()
}
