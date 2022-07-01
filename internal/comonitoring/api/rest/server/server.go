package server

import (
	"comonitoring/config"
	"comonitoring/internal/comonitoring/metrics"
	middleware "comonitoring/internal/middlewares"
	"comonitoring/pkg/logger"
	"fmt"
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"net/http"
)

const ApiV1 = "/api/v1"

type Module interface {
	RegisterRoutes(group *gin.RouterGroup)
}

var log, _ = logger.NewLogger(logger.Config{})

func NewServer(healthModule, metricModule Module, cfg config.Config) (*http.Server, error) {
	engine := gin.Default()

	// CORS
	engine.Use(middleware.NewCORSMiddleware().Middleware())
	log.Debugln("The CORS middleware has been set")

	// Register system metrics
	err := metrics.RegisterAll()
	if err != nil {
		return nil, err
	}

	// Prometheus Exporter
	promExporter := ginprom.New(
		ginprom.Engine(engine),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	engine.Use(promExporter.Instrument())

	// Register the handlers
	v1 := engine.Group(ApiV1)
	healthModule.RegisterRoutes(v1)
	metricModule.RegisterRoutes(v1)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.COMonitoring.ServingPort),
		Handler: engine,
	}
	log.Infoln("new server has been created")

	return httpServer, nil
}
