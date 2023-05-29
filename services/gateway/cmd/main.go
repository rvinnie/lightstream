package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rvinnie/lightstream/services/gateway/transport/rest"
	"github.com/rvinnie/lightstream/services/gateway/transport/rest/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	gin.SetMode(cfg.GIN.Mode)

	// Initializing gRPC connection
	grpcTarget := fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.Port)
	grpcConn, err := grpc.Dial(grpcTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatal(err)
	}
	defer grpcConn.Close()

	handlers := handler.NewHandler(grpcConn)
	restServer := rest.NewServer(cfg, handlers.InitRoutes(*cfg))

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

	logrus.Info("Gateway (HTTP) server shutting down")

	if err = restServer.Stop(context.Background()); err != nil {
		logrus.Errorf("Error on gateway (HTTP) server shutting down: %s", err.Error())
	}
}
