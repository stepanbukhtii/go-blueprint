package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/stepanbukhtii/easy-tools/config"
)

type Config struct {
	Service       config.Service
	Log           config.Log
	Database      config.DB
	Redis         config.Redis
	OpenTelemetry config.OpenTelemetry
	API           config.API
	JWT           config.JWT
	GRPC          config.GRPC
	RabbitMQ      config.RabbitMQ
	Kafka         config.Kafka
	NATS          config.NATS
	RandomUser    RandomUser
}

type RandomUser struct {
	BaseURL string `env:"RANDOM_USER_BASE_URL"`
}

func Read() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
