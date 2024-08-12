package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PGConfig
	NatsConfig
	ServerConfig
	RedisConfig
}

type PGConfig struct {
	PGName     string
	PGPassword string
	PGUser     string
	PGHost     string
	PGPort     string
}

type NatsConfig struct {
	ClusterId string
	ClientID  string
	Url       string
}

type ServerConfig struct {
	HTTPServerPort string
	ServerAddress  string
	Env            string
}

type RedisConfig struct {
	RedisDB       string
	RedisPassword string
	RedisPort     string
	RedisDuration string
}

func MustLoad() *Config {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	return &Config{
		PGConfig: PGConfig{
			PGName:     os.Getenv("PG_DB_NAME"),
			PGPassword: os.Getenv("PG_PASSWORD"),
			PGUser:     os.Getenv("PG_USER"),
			PGHost:     os.Getenv("PG_HOST"),
			PGPort:     os.Getenv("PG_PORT"),
		},
		NatsConfig: NatsConfig{
			ClusterId: os.Getenv("NATS_CLUSTER_ID"),
			ClientID:  os.Getenv("NATS_CLIENT_ID"),
			Url:       os.Getenv("NATS_URL"),
		},
		ServerConfig: ServerConfig{
			HTTPServerPort: os.Getenv("HTTP_SERVER_PORT"),
			ServerAddress:  os.Getenv("HTTP_SERVER_ADDRESS"),
			Env:            os.Getenv("ENV"),
		},
		RedisConfig: RedisConfig{
			RedisDB:       os.Getenv("REDIS_DB"),
			RedisPassword: os.Getenv("REDIS_PASSWORD"),
			RedisPort:     os.Getenv("REDIS_PORT"),
			RedisDuration: os.Getenv("REDIS_DURATION"),
		},
	}
}
