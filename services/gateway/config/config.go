package config

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	HTTP HTTPConfig
	GRPC GRPCConfig
	GIN  GINConfig
}

type HTTPConfig struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type GRPCConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type GINConfig struct {
	Mode string
}

func InitConfig(configDir string) (*Config, error) {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("gateway")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("gRPC", &cfg.GRPC); err != nil {
		return nil, err
	}

	setEnvVariables(&cfg)

	return &cfg, nil
}

func setEnvVariables(cfg *Config) {
	cfg.GIN.Mode = os.Getenv("GIN_MODE")
}
