package storage

import (
	"github.com/gin-gonic/gin"
	"github.com/rvinnie/lightstream/services/storage/config"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	awsManager *AWSManager
}

func NewHandler(awsManager *AWSManager) *Handler {
	return &Handler{awsManager: awsManager}
}

func (h *Handler) InitRoutes(cfg config.Config) *gin.Engine {
	gin.SetMode(cfg.GIN.Mode)
	router := gin.New()

	router.GET("/image", h.image)

	return router
}

func (h *Handler) image(c *gin.Context) {
	const path = "images/eg.jpg"
	// TODO: add validation of path

	object, err := h.awsManager.DownloadObject(path)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.Data(http.StatusOK, object.contentType, object.body)
}
