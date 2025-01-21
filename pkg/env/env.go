package env

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

func New[T any](cfg T, filenames ...string) (T, error) {
	err := godotenv.Load(filenames...)
	if err != nil && !os.IsNotExist(err) {
		return cfg, err
	}

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
