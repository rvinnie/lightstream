package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rvinnie/lightstream/services/gateway/config"
	pb "github.com/rvinnie/lightstream/services/gateway/pb"
	"github.com/rvinnie/lightstream/services/gateway/service"
	"google.golang.org/grpc"
)

type ImagesHandler struct {
	imageStorageClient pb.ImageStorageClient
	imagesService      *service.ImagesService
}

func NewImagesHandler(grpcConn grpc.ClientConnInterface, imagesService *service.ImagesService) *ImagesHandler {
	return &ImagesHandler{
		imageStorageClient: pb.NewImageStorageClient(grpcConn),
		imagesService:      imagesService,
	}
}

func (h *ImagesHandler) InitRoutes(cfg config.Config) *gin.Engine {
	gin.SetMode(cfg.GIN.Mode)
	router := gin.New()

	router.GET("/image/:path", h.image)

	return router
}

func (h *ImagesHandler) image(c *gin.Context) {
	param := c.Param("path")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid path"))
		return
	}

	path, err := h.imagesService.GetById(c, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	resp, err := h.imageStorageClient.GetImage(c, &pb.ImageStorageRequest{Path: path})

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, resp.ContentType, resp.Image)
}
