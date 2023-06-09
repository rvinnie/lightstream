package handler

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rvinnie/lightstream/services/gateway/transport/amqp"
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
	rabbitProducer     *amqp.Producer
}

func NewImagesHandler(grpcConn grpc.ClientConnInterface, imagesService *service.ImagesService, rabbitProducer *amqp.Producer) *ImagesHandler {
	return &ImagesHandler{
		imageStorageClient: pb.NewImageStorageClient(grpcConn),
		imagesService:      imagesService,
		rabbitProducer:     rabbitProducer,
	}
}

func (h *ImagesHandler) InitRoutes(cfg config.Config) *gin.Engine {
	gin.SetMode(cfg.GIN.Mode)
	router := gin.New()

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
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

	err = h.rabbitProducer.Publish("msg from producer: " + resp.ContentType)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.Data(http.StatusOK, resp.ContentType, resp.Image)
}
