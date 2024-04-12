package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func Load(filenames ...string) error {
	if err := godotenv.Load(filenames...); err != nil {
		return fmt.Errorf("can't load env files: %w", err)
	}
	return nil
}

func Parse[T any](prefix string) (T, error) {
	var conf T

	if err := envconfig.Process(prefix, &conf); err != nil {
		return *new(T), err
	}

	return conf, nil
}
