package rest

import (
	"comonitoring/config"
	"comonitoring/internal/comonitoring/metrics"
	middleware "comonitoring/internal/middlewares"
	"comonitoring/pkg/logger"
	"fmt"
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

var log, _ = logger.NewLogger(logger.Config{})

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type CODensityRequest struct {
	Density *float64 `json:"density"`
}

func PushCODensity(c *gin.Context) {
	// Parse params
	orgUUIDParam := c.Param("orgUUID")
	orgUUID, err := uuid.FromString(orgUUIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not push this value",
			"error":   err.Error(),
		})
		return
	}
	deviceUUIDParam := c.Param("deviceUUID")
	deviceUUID, err := uuid.FromString(deviceUUIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not push this value",
			"error":   err.Error(),
		})
		return
	}

	// Parse request
	CODensity := CODensityRequest{}
	err = c.BindJSON(&CODensity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not push this value",
			"error":   err.Error(),
		})
		return
	}

	// Store it in server metrics
	if CODensity.Density == nil {
		log.Debugf("Recieved a no data CO Density value from org: %s device: %s", orgUUIDParam, deviceUUIDParam)
		metrics.COGauge.DeleteLabelValues(orgUUID.String(), deviceUUID.String())

		c.JSON(http.StatusOK, gin.H{
			"message": "No data pushed successfully",
		})
	} else {
		log.Debugf("Recieved a new CO Density value: %f from org: %s device: %s", *CODensity.Density, orgUUIDParam, deviceUUIDParam)
		metrics.COGauge.WithLabelValues(orgUUID.String(), deviceUUID.String()).Set(*CODensity.Density)

		c.JSON(http.StatusOK, gin.H{
			"message": "CO Density pushed successfully",
		})
	}
}
func NewServer(cfg config.Config) (*http.Server, error) {
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
	v1 := engine.Group("/api/v1")
	v1.GET("/ping", HealthCheck)
	v1.PUT("/organizations/:orgUUID/devices/:deviceUUID/metrics/co_density", PushCODensity)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.COMonitoring.ServingPort),
		Handler: engine,
	}

	return httpServer, nil
}
