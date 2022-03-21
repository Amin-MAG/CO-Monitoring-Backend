package rest

import (
	"comonitoring/config"
	middleware "comonitoring/internal/middlewares"
	"comonitoring/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var log, _ = logger.NewLogger(logger.Config{})

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func NewServer(cfg config.Config) (*http.Server, error) {
	engine := gin.Default()

	// CORS
	engine.Use(middleware.NewCORSMiddleware().Middleware())
	log.Debugln("The CORS middleware has been set")

	// Register the handlers
	v1 := engine.Group("/api/v1")
	v1.GET("/ping", HealthCheck)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.COMonitoring.ServingPort),
		Handler: engine,
	}

	return httpServer, nil
}
