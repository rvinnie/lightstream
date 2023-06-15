package handler

import (
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rvinnie/lightstream/services/gateway/monitoring"
	"github.com/rvinnie/lightstream/services/gateway/transport/amqp"
	"io/ioutil"
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

const (
	imagesDirectoryName = "images"
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
		metrics.ConcurrentRequests.Inc()
		method := c.Request.Method
		url := c.Request.URL.Path

		// Truncate resource ids
		if unicode.IsDigit(rune(url[len(url)-1])) {
			url = filepath.Dir(url)
		}

		c.Next()

		metrics.ConcurrentRequests.Dec()
		statusCode := c.Writer.Status()
		duration := time.Since(start).Seconds()
		metrics.CollectMetrics(method, url, statusCode, duration)
	}
}

func (h *ImagesHandler) InitRoutes(cfg config.Config) *gin.Engine {
	gin.SetMode(cfg.GIN.Mode)
	router := gin.New()
	router.Use(PrometheusMiddleware(h.metrics))
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type", "Accept-Encoding", "Filename"},
	}))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/images/:path", h.getImage)
	router.GET("/images", h.getImages)
	router.POST("/images/add", h.createImage)

	return router
}

type imageResponse struct {
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	Data        []byte `json:"data"`
}

func (h *ImagesHandler) createImage(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	filename := c.Request.Header.Get("Filename")
	if filename == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	path := imagesDirectoryName + "/" + filename
	id, err := h.imagesService.Create(c, path)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = h.imageStorageClient.CreateImage(c, &pb.CreateImageRequest{
		Path:        path,
		ContentType: contentType,
		Image:       data,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, id)
}

func (h *ImagesHandler) getImage(c *gin.Context) {
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

	resp, err := h.imageStorageClient.GetImage(c, &pb.FindImageRequest{Path: path})
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = h.rabbitProducer.Publish(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//c.JSON(http.StatusOK, imageResponse{
	//	Name:        resp.Name,
	//	ContentType: resp.ContentType,
	//	Data:        resp.Image,
	//})
	c.Data(http.StatusOK, resp.ContentType, resp.Image)
}

func (h *ImagesHandler) getImages(c *gin.Context) {
	paths, err := h.imagesService.GetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resp, err := h.imageStorageClient.GetImages(c, &pb.FindImagesRequest{Paths: paths})
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var images []imageResponse
	for _, image := range resp.GetImages() {
		images = append(images, imageResponse{
			Name:        image.Name,
			ContentType: image.ContentType,
			Data:        image.Image,
		})
	}

	c.JSON(http.StatusOK, images)
}
