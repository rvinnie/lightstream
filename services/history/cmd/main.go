package main

import (
	postgres "github.com/rvinnie/lightstream/pkg/database"
	"github.com/rvinnie/lightstream/services/history/transport/amqp"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rvinnie/lightstream/services/history/config"
	"github.com/sirupsen/logrus"
)

const (
	configPath = "./config"
)

func main() {
	// Adding logger
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Initializing env variables
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	// Initializing config
	cfg, err := config.InitConfig(configPath)
	if err != nil {
		logrus.Fatal("Unable to parse config", err)
	}

	// Initializing postgres
	db, err := postgres.NewConnPool(postgres.DBConfig{
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		DBName:   cfg.Postgres.DBName,
	})
	if err != nil {
		logrus.Errorf("Unable to connect db: %v", err)
		return
	}
	defer db.Close()

	// Initializing RabbitMQ consumer
	consumer, err := amqp.NewConsumer(amqp.ConsumerConfig{
		Username: cfg.RabbitMQ.Username,
		Password: cfg.RabbitMQ.Password,
		Host:     cfg.RabbitMQ.Host,
		Port:     cfg.RabbitMQ.Port,
	})
	if err != nil {
		logrus.Errorf("Unable to create RabbitMQ consumer: %v", err)
		return
	}

	if err = consumer.Consume(); err != nil {
		logrus.Fatal("Consuming failed: ", err)
	}
	logrus.Info("History (RabbitMQ) consumer is running")

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	<-quit

	logrus.Info("History (RabbitMQ) consumer shutting down")
	if err = consumer.Shutdown(); err != nil {
		logrus.Errorf("Error on history (RabbitMQ) consumer shutting down: %s", err.Error())
	}
}
