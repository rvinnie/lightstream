package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/rvinnie/lightstream/services/storage/aws"
	"github.com/rvinnie/lightstream/services/storage/config"
	pb "github.com/rvinnie/lightstream/services/storage/pb"
)

type ImageStorageService interface {
	CreateImage(ctx context.Context, request *pb.CreateImageRequest) (*emptypb.Empty, error)
	GetImage(ctx context.Context, request *pb.FindImageRequest) (*pb.FindImageResponse, error)
	GetImages(ctx context.Context, request *pb.FindImagesRequest) (*pb.FindImagesResponse, error)
}

type ImageStorageHandler struct {
	pb.UnimplementedImageStorageServer

	manager aws.AWS
	cfg     config.AWSConfig
}

func NewImageStorageHandler(m aws.AWS, cfg config.AWSConfig) *ImageStorageHandler {
	return &ImageStorageHandler{manager: m, cfg: cfg}
}

func (h *ImageStorageHandler) CreateImage(ctx context.Context, request *pb.CreateImageRequest) (*emptypb.Empty, error) {
	err := h.manager.UploadObject(request.Path, request.ContentType, request.Image)

	return &emptypb.Empty{}, err
}

func (h *ImageStorageHandler) GetImage(ctx context.Context, request *pb.FindImageRequest) (*pb.FindImageResponse, error) {
	object, err := h.manager.DownloadObject(request.Path)

	imageResponse := &pb.FindImageResponse{
		Image:       object.Body,
		ContentType: object.ContentType,
		Name:        object.Name,
	}

	return imageResponse, err
}

func (h *ImageStorageHandler) GetImages(ctx context.Context, request *pb.FindImagesRequest) (*pb.FindImagesResponse, error) {
	objects, err := h.manager.DownloadObjects(request.Paths)

	var imagesResponses []*pb.FindImageResponse
	for _, object := range objects {
		imageResponse := &pb.FindImageResponse{
			Name:        object.Name,
			ContentType: object.ContentType,
			Image:       object.Body,
		}
		imagesResponses = append(imagesResponses, imageResponse)
	}

	return &pb.FindImagesResponse{Images: imagesResponses}, err
}
