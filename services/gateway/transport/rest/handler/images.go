package handler

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rvinnie/lightstream/services/gateway/monitoring"
	"github.com/rvinnie/lightstream/services/gateway/transport/amqp"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
	"unicode"

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
	metrics            *monitoring.Metrics
}

func NewImagesHandler(grpcConn grpc.ClientConnInterface, imagesService *service.ImagesService, rabbitProducer *amqp.Producer, metrics *monitoring.Metrics) *ImagesHandler {
	return &ImagesHandler{
		imageStorageClient: pb.NewImageStorageClient(grpcConn),
		imagesService:      imagesService,
		rabbitProducer:     rabbitProducer,
		metrics:            metrics,
	}
}

func PrometheusMiddleware(metrics *monitoring.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		url := c.Request.URL.Path

		// Truncate resource ids
		if unicode.IsDigit(rune(url[len(url)-1])) {
			url = filepath.Dir(url)
		}

		c.Next()

		statusCode := c.Writer.Status()
		duration := time.Since(start).Seconds()
		metrics.CollectMetrics(method, url, statusCode, duration)
	}
}

func (h *ImagesHandler) InitRoutes(cfg config.Config) *gin.Engine {
	gin.SetMode(cfg.GIN.Mode)
	router := gin.New()
	router.Use(PrometheusMiddleware(h.metrics))

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
		return
	}

	c.Data(http.StatusOK, resp.ContentType, resp.Image)
}
