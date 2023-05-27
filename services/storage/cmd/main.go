package main

import (
	"github.com/rvinnie/lightstream/services/storage/aws"
	"github.com/rvinnie/lightstream/services/storage/config"
	"github.com/rvinnie/lightstream/services/storage/transport/grpc"
	"github.com/rvinnie/lightstream/services/storage/transport/grpc/handler"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
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

	//Initializing config
	cfg, err := config.InitConfig(configPath)
	if err != nil {
		logrus.Error("Unable to parse config", err)
		return
	}

	// Creating AWS manager
	awsManager := aws.NewAWSManager(cfg.AWS.BucketName, cfg.AWS.Config)

	// Creating handlers
	grpcHandler := handler.NewImageStorageHandler(*awsManager, cfg.AWS)

	// Creating gRPC server
	grpcServer := grpc.NewServer(grpcHandler)
	go func() {
		logrus.Info("Starting storage microservice")
		if err = grpcServer.ListenAndServe(cfg.GRPC.Port); err != nil {
			logrus.Fatalf("error occured while running grpc server: %s", err.Error())
		}
	}()
	logrus.Info("Storage microservice is running")

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	<-quit
	logrus.Info("Storage microservice shutting down")

	grpcServer.Stop()
}
