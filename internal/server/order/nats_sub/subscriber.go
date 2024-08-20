package nats_sub

import (
	"L0/entities"
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log/slog"
	"time"
)

type RepositoryGiver interface {
	GetOrderById(con context.Context, orderUID string) (*entities.Order, error)
	CreateOrder(ctx context.Context, order *entities.Order) error
}

const (
	retryAttempts = 3
	retryDelay    = 1 * time.Second
)

type OrderSubscriber struct {
	stanConn stan.Conn
	log      *slog.Logger
	repos    RepositoryGiver
}

func NewOrderSubscriber(stanConn stan.Conn, log *slog.Logger, repos RepositoryGiver) *OrderSubscriber {
	return &OrderSubscriber{
		stanConn: stanConn,
		log:      log,
		repos:    repos,
	}
}

func (s *OrderSubscriber) Subscribe(subject, qgroup string) {
	s.log.Info("Subscribing to subject", "subject", subject, "qgroup", qgroup)

	s.stanConn.Subscribe(subject, func(m *stan.Msg) {
		var order entities.Order
		err := json.Unmarshal(m.Data, &order)
		if err != nil {
			s.log.Warn("Error unmarshalling message", "error", err)
			return
		}
		err = s.repos.CreateOrder(context.Background(), &order)
		if err != nil {
			s.log.Warn("Error creating order", "error", err)
		}
	})

	s.log.Info("Message recieved succeffully", "subject", subject, "qgroup", qgroup)
}
