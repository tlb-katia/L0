package server

import (
	"L0/config"
	"L0/entities"
	nats_sub2 "L0/internal/server/order/nats_sub"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log/slog"
	"net/http"
)

type RepositoryGenerator interface {
	GetOrderById(con context.Context, orderUID string) (*entities.Order, error)
	CreateOrder(ctx context.Context, order *entities.Order) error
}

type NATSGenerator interface {
	MessageHandler(msg *nats.Msg)
	Subscribe() (*nats.Subscription, error)
	ValidateOrderData(order *entities.Order) bool
}

type Server struct {
	db     RepositoryGenerator
	stan   stan.Conn
	router *chi.Mux
	log    *slog.Logger
}

func NewServer(db RepositoryGenerator, stan stan.Conn, router *chi.Mux, log *slog.Logger) *Server {
	return &Server{
		db:     db,
		stan:   stan,
		router: router,
		log:    log,
	}
}

func (s *Server) Run(config *config.Config) {
	log := s.log.With(
		slog.String("method", "serverRun"))

	go func() {
		sb := nats_sub2.NewOrderSubscriber(s.stan, log, s.db)
		sb.Subscribe(nats_sub2.CreateOrderSubject, nats_sub2.OrderGroupName)
	}()

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.URLFormat)

	s.router.Get("/{id}", s.GetOrderById)
	s.router.Get("/", s.SayHello)

	srv := http.Server{
		Addr:    config.HTTPServerPort,
		Handler: s.router,
	}

	if err := srv.ListenAndServe(); err != nil {
		s.log.Error("failed to start server")
	}
}
