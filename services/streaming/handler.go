package streaming

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rvinnie/lightstream/services/streaming/config"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) InitRoutes(cfg config.Config) *gin.Engine {
	gin.SetMode(cfg.GIN.Mode)
	router := gin.New()

	router.GET("/video", h.video)

	return router
}

func (h *Handler) video(c *gin.Context) {
	const path = "./videos/test.mp4"

	_, err := os.ReadFile(path)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.File(path)
}
