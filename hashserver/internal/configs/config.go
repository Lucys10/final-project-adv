package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	GrpcServerAddress string `env:"GRPC_SERVER_ADDRESS"`
	Network           string `env:"NETWORK"`
}

var config Config

func GetConfig() (*Config, error) {
	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, errors.WithStack(err)
	}
	return &config, nil
}
