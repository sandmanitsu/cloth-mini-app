package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host  string `env:"HOST" env-required:"true"`
	Port  string `env:"PORT" env-required:"true"`
	Env   string `env:"ENV" env-required:"true"`
	DB    DB
	Minio Minio
}

type DB struct {
	Host     string `env:"DBHOST" env-required:"true"`
	User     string `env:"USER" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	Port     int    `env:"DBPORT" env-required:"true"`
	DBname   string `env:"DBNAME" env-required:"true"`
}

type Minio struct {
	Endpoint   string `env:"MINIO_ENDPOINT" env-required:"true"`
	BucketName string `env:"MINIO_BUCKET_NAME" env-required:"true"`
	User       string `env:"MINIO_ROOT_USER" env-required:"true"`
	Password   string `env:"MINIO_ROOT_PASSWORD" env-required:"true"`
}

var (
	config *Config
	once   sync.Once
)

// Getting config variables from .env file
func MustLoad(path string) *Config {
	if config == nil {
		once.Do(
			func() {
				configPath := filepath.Join(path)

				if _, err := os.Stat(configPath); err != nil {
					log.Fatalf("Error opening config file: %s", err)
				}

				var newConfig Config
				err := cleanenv.ReadConfig(configPath, &newConfig)
				if err != nil {
					log.Fatalf("Error reading config file: %s", err)
				}

				config = &newConfig
			})
	}

	return config
}
