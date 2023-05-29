package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rvinnie/lightstream/services/gateway/config"
	pb "github.com/rvinnie/lightstream/services/gateway/pb"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

const (
	imageStorageFolder = "images"
)

type Handler struct {
	imageStorageClient pb.ImageStorageClient
}

func NewHandler(grpcConn grpc.ClientConnInterface) *Handler {
	return &Handler{imageStorageClient: pb.NewImageStorageClient(grpcConn)}
}

func (h *Handler) InitRoutes(cfg config.Config) *gin.Engine {
	gin.SetMode(cfg.GIN.Mode)
	router := gin.New()

	router.GET("/image/:path", h.image)

	return router
}

func (h *Handler) image(c *gin.Context) {
	param := c.Param("path")
	if !isImage(param) {
		c.AbortWithError(http.StatusBadRequest, errors.New("file is not an image (unavailable image extension)"))
		return
	}

	path := fmt.Sprintf("%s/%s", imageStorageFolder, param)
	resp, err := h.imageStorageClient.GetImage(c, &pb.ImageStorageRequest{Path: path})

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, resp.ContentType, resp.Image)
}

func isImage(param string) bool {
	exts := []string{
		"gif", "jpeg", "jpg", "pjpeg", "jfif", "pjp", "png", "apng",
		"avif", "svg", "webp", "tiff", "tif", "bmp", "ico", "cur",
	}

	splittedParam := strings.Split(param, ".")
	countOfPieces := len(splittedParam)
	ext := splittedParam[countOfPieces-1]

	if countOfPieces < 2 {
		return false
	}

	for _, curExt := range exts {
		if curExt == ext {
			return true
		}
	}
	return false
}
