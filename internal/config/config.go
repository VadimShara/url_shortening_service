package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string `env:"ENV" env-default:"local"`
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`

	GRPC GRPCConfig

	Storage string `env:"STORAGE" env-default:""`

	PGConfig       PostresConfig
	MigrationsPath string `env:"MIGRATIONS_PATH" env-default:"migrations"`
}

type PostresConfig struct {
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     int    `env:"POSTGRES_PORT" env-required:"true"`
	DBName   string `env:"POSTGRES_DBNAME" env-required:"true"`
	SSLMode  string `env:"POSTGRES_SSLMODE" env-default:"disable"`
}

type GRPCConfig struct {
	Port    int           `env:"PORT" env-required:"true"`
	Timeout time.Duration `env:"TIMEOUT" env-default:"10h"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func Load() Config {
	var config Config

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatal("couldn't bind settings to config")
	}

	return config
}
