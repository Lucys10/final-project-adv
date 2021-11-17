package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	GrpcServerAddress string `env:"GRPC_SERVER_ADDRESS"`
	MongoURL          string `env:"MONGO_URL"`
	ServerAddress     string `env:"SERVER_ADDRESS"`
}

var config Config

func GetConfig() (*Config, error) {

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, errors.WithStack(err)
	}

	return &config, nil
}
