package handler

import (
	"context"
	"github.com/rvinnie/lightstream/services/storage/aws"
	"github.com/rvinnie/lightstream/services/storage/config"
	pb "github.com/rvinnie/lightstream/services/storage/pb"
)

type ImageStorageService interface {
	GetImage(ctx context.Context, request *pb.ImageStorageRequest) (*pb.ImageStorageResponse, error)
}

type ImageStorageHandler struct {
	pb.UnimplementedImageStorageServer

	manager aws.AWSManager
	cfg     config.AWSConfig
}

func NewImageStorageHandler(m aws.AWSManager, cfg config.AWSConfig) *ImageStorageHandler {
	return &ImageStorageHandler{manager: m, cfg: cfg}
}

func (h *ImageStorageHandler) GetImage(ctx context.Context, request *pb.ImageStorageRequest) (*pb.ImageStorageResponse, error) {
	awsManager := aws.NewAWSManager(h.cfg.BucketName, h.cfg.Config)
	object, err := awsManager.DownloadObject(request.Path)

	return &pb.ImageStorageResponse{Image: object.Body, ContentType: object.ContentType}, err
}
