package metric

import (
	"comonitoring/internal/comonitoring/api/rest/server"
	"comonitoring/internal/comonitoring/metrics"
	absms "comonitoring/internal/sms"
	"comonitoring/pkg/logger"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

var log, _ = logger.NewLogger(logger.Config{})

type Metrics struct {
	smsProvider *absms.SMS
}

func NewModule(smsProvider *absms.SMS) (server.Module, error) {
	return &Metrics{
		smsProvider: smsProvider,
	}, nil
}

func (m *Metrics) PushCODensity(c *gin.Context) {
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

func (m *Metrics) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/organizations/:orgUUID/devices/:deviceUUID/metric/co_density", m.PushCODensity)
}
