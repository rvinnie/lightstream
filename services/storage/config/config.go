package config

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	GRPC GRPCConfig
	AWS  AWSConfig
	GIN  GINConfig
}

type GRPCConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AWSConfig struct {
	BucketName string `yaml:"bucketName"`
	Config     aws.Config
}

type GINConfig struct {
	Mode string
}

func InitConfig(configDir string) (*Config, error) {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("storage")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.UnmarshalKey("gRPC", &cfg.GRPC); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("aws", &cfg.AWS); err != nil {
		return nil, err
	}

	loadAWSConfig(&cfg)
	setEnvVariables(&cfg)

	return &cfg, nil
}

func loadAWSConfig(cfg *Config) error {
	// Create a custom endpoint resolver
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID && region == "ru-central1" {
			return aws.Endpoint{
				PartitionID:   "yc",
				URL:           "https://storage.yandexcloud.net",
				SigningRegion: "ru-central1",
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	// Load config from ~/.aws/*
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		return err
	}

	cfg.AWS.Config = awsCfg

	return err
}

func setEnvVariables(cfg *Config) {
	cfg.GIN.Mode = os.Getenv("GIN_MODE")
}
