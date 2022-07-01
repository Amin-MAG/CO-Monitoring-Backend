package health

import (
	"comonitoring/internal/comonitoring/api/rest/server"
	"github.com/gin-gonic/gin"
)

type Health struct{}

func NewModule() (server.Module, error) {
	return &Health{}, nil
}

func (h *Health) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *Health) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/ping", h.HealthCheck)
}
