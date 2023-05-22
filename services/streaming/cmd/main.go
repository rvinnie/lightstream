package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/rvinnie/lightstream/services/streaming"
	"github.com/rvinnie/lightstream/services/streaming/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	configPath = "services/streaming/config"
	envPath    = "services/streaming/.env"
)

func main() {
	// Adding logger
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Initializing env variables
	if err := godotenv.Load(envPath); err != nil {
		logrus.Error("Error loading .env file")
		return
	}

	// Initializing config
	cfg, err := config.InitConfig(configPath)
	if err != nil {
		logrus.Error("Unable to parse config", err)
		return
	}

	handler := streaming.NewHandler()
	server := &http.Server{
		Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:      handler.InitRoutes(*cfg),
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	go func() {
		server.ListenAndServe()
	}()
	logrus.Info("Video streaming microservice is running")

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	<-quit
	logrus.Info("Video streaming microservice shutting down")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error on streaming microservice shutting down: %s", err.Error())
	}
}
