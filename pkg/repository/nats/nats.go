package nats

import (
	"L0/config"
	"github.com/nats-io/stan.go"
	"log/slog"
	"time"
)

const (
	connectWait        = time.Second * 30
	pubAckWait         = time.Second * 30
	interval           = 10
	maxOut             = 5
	maxPubAcksInflight = 25
)

func NewNatsConnect(cfg *config.Config, log *slog.Logger) (stan.Conn, error) {
	return stan.Connect(
		cfg.NatsConfig.ClusterId,
		cfg.NatsConfig.ClientID,
		stan.ConnectWait(connectWait),
		stan.PubAckWait(pubAckWait),
		stan.NatsURL(cfg.NatsConfig.Url),
		stan.Pings(interval, maxOut),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			log.Error("Connection lost, reason: %v", err)
		}),
		stan.MaxPubAcksInflight(maxPubAcksInflight),
	)
}
