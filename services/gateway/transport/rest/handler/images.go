package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rvinnie/lightstream/services/gateway/config"
	pb "github.com/rvinnie/lightstream/services/gateway/pb"
	"github.com/rvinnie/lightstream/services/gateway/service"
	"google.golang.org/grpc"
)

const (
	imageStorageFolder = "images"
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
	// if !isImage(param) {
	// 	c.AbortWithError(http.StatusBadRequest, errors.New("file is not an image (unavailable image extension)"))
	// 	return
	// }

	//path := fmt.Sprintf("%s/%s", imageStorageFolder, param)
	//path := h.
	path, err := h.imagesService.GetById(c, param)
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

// func isImage(param string) bool {
// 	exts := []string{
// 		"gif", "jpeg", "jpg", "pjpeg", "jfif", "pjp", "png", "apng",
// 		"avif", "svg", "webp", "tiff", "tif", "bmp", "ico", "cur",
// 	}

// 	splittedParam := strings.Split(param, ".")
// 	countOfPieces := len(splittedParam)
// 	ext := splittedParam[countOfPieces-1]

// 	if countOfPieces < 2 {
// 		return false
// 	}

// 	for _, curExt := range exts {
// 		if curExt == ext {
// 			return true
// 		}
// 	}
// 	return false
// }
