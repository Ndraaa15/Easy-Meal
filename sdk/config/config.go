package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Initializer interface {
	LoadEnv(envPath string) error
	GetEnv(key string) string
}

type config struct{}

func Init() Initializer {
	return &config{}
}

func (c *config) LoadEnv(envPath string) error {
	return godotenv.Load(envPath)
}

func (c *config) GetEnv(key string) string {
	return os.Getenv(key)
}
