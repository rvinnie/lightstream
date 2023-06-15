package main

import (
	"context"
	"fmt"
	postgres "github.com/rvinnie/lightstream/pkg/database"
	"github.com/rvinnie/lightstream/services/gateway/monitoring"
	"github.com/rvinnie/lightstream/services/gateway/transport/amqp"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rvinnie/lightstream/services/gateway/repository"
	"github.com/rvinnie/lightstream/services/gateway/service"
	"github.com/rvinnie/lightstream/services/gateway/transport/rest"
	"github.com/rvinnie/lightstream/services/gateway/transport/rest/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/joho/godotenv"

	"github.com/rvinnie/lightstream/services/gateway/config"
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

	// Initializing RabbitMQ producer
	rabbitProducer, err := amqp.NewProducer(amqp.ProducerConfig{
		Username: cfg.RabbitMQ.Username,
		Password: cfg.RabbitMQ.Password,
		Host:     cfg.RabbitMQ.Host,
		Port:     cfg.RabbitMQ.Port,
	})
	if err != nil {
		logrus.Errorf("Unable to create RabbitMQ producer: %v", err)
		return
	}
	logrus.Info("History RabbitMQ (AMQP) producer is created")

	// Initializing gRPC connection
	grpcTarget := fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.Port)
	grpcConn, err := grpc.Dial(grpcTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatal(err)
	}
	defer grpcConn.Close()
	logrus.Info("Storage (gRPC) client is created")

	// Initializing Prometheus
	metrics := monitoring.CreateMetrics(cfg.Prometheus.MaxConcurrentRequests)

	imagesRepository := repository.NewImagesPostgres(db)
	imagesService := service.NewImagesService(imagesRepository)
	imagesHandler := handler.NewImagesHandler(grpcConn, imagesService, rabbitProducer, metrics)

	restServer := rest.NewServer(cfg, imagesHandler.InitRoutes(*cfg))
	go func() {
		if err = restServer.Run(); err != http.ErrServerClosed {
			logrus.Fatalf("error occured while running gateway (HTTP) server: %s", err.Error())
		}
	}()
	logrus.Info("Gateway (HTTP) server is running")

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	<-quit

	logrus.Info("History RabbitMQ (AMQP) producer shutting down")
	if err = rabbitProducer.Shutdown(); err != nil {
		logrus.Errorf("Error on history RabbitMQ (AMQP) producer shutting down: %s", err.Error())
	}

	logrus.Info("Gateway (HTTP) server shutting down")
	if err = restServer.Stop(context.Background()); err != nil {
		logrus.Errorf("Error on gateway (HTTP) server shutting down: %s", err.Error())
	}
}
