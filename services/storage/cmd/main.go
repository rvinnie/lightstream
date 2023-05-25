package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rvinnie/lightstream/services/storage"
	"github.com/rvinnie/lightstream/services/storage/config"
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
		logrus.Error("Error loading .env file")
		return
	}

	// Initializing config
	cfg, err := config.InitConfig(configPath)
	if err != nil {
		logrus.Error("Unable to parse config", err)
		return
	}

	// Creating AWS manager
	awsManager := storage.NewAWSManager(cfg.AWS.BucketName, cfg.AWS.Config)

	handler := storage.NewHandler(awsManager)
	server := &http.Server{
		Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:      handler.InitRoutes(*cfg),
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Info("Storage microservice is running")

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	<-quit
	logrus.Info("Storage microservice shutting down")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error on storage microservice shutting down: %s", err.Error())
	}
}
