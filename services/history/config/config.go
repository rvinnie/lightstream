package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Postgres PostgresConfig
	RabbitMQ RabbitMQConfig
}

type PostgresConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

type RabbitMQConfig struct {
	Username string
	Password string
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

func InitConfig(configDir string) (*Config, error) {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("history")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.UnmarshalKey("rabbit", &cfg.RabbitMQ); err != nil {
		return nil, err
	}

	setEnvVariables(&cfg)

	return &cfg, nil
}

func setEnvVariables(cfg *Config) {
	cfg.Postgres.Username = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.Host = os.Getenv("DATABASE_HOST")
	cfg.Postgres.DBName = os.Getenv("POSTGRES_DB")

	cfg.RabbitMQ.Username = os.Getenv("RABBIT_USER")
	cfg.RabbitMQ.Password = os.Getenv("RABBIT_PASSWORD")
}
