package main

import (
	"L0/config"
	"L0/pkg/repository/db"
	"L0/pkg/repository/nats"
	"L0/pkg/repository/redis"
	"L0/server"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	conf := config.MustLoad()
	log := setupLogger(conf.Env)
	log.Info("starting server")

	psql, err := db.Init(conf)
	if err != nil {
		log.Error("Failed to init database", err)
		os.Exit(1)
	}
	defer psql.Close()

	rdb, err := redis.NewRedis(conf)
	if err != nil {
		log.Error("failed to init redis", err)
		os.Exit(1)
	}

	natsConn, err := nats.NewNatsConnect(conf, log)
	if err != nil {
		log.Error("failed to init nats_sub connection", err)
		os.Exit(1)
	}

	repo := db.NewRepository(psql, rdb.Client)
	srv := server.NewServer(repo, natsConn, chi.NewRouter(), log)
	srv.Run(conf)

	log.Info("server stopped")

	//fmt.Println("I can work!")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
