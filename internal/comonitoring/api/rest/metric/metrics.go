package metric

import (
	"comonitoring/internal/comonitoring/api/rest/server"
	"comonitoring/internal/comonitoring/metrics"
	absms "comonitoring/internal/sms"
	"comonitoring/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

var log, _ = logger.NewLogger(logger.Config{})

type Metrics struct {
	smsProvider absms.SMS
	smsCache    map[string]time.Time
}

func NewModule(smsProvider absms.SMS, smsCache map[string]time.Time) (server.Module, error) {
	return &Metrics{
		smsProvider: smsProvider,
		smsCache:    smsCache,
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

	// Handle empty density
	if CODensity.Density == nil {
		log.Debugf("Recieved a no data CO Density value from org: %s device: %s", orgUUIDParam, deviceUUIDParam)
		metrics.COGauge.DeleteLabelValues(orgUUID.String(), deviceUUID.String())

		c.JSON(http.StatusOK, gin.H{
			"message": "No data pushed successfully",
		})
		return
	}

	// Check alerting
	if *CODensity.Density < 1.0 {
		log.Infof("%+v", m.smsCache)

		now := time.Now()
		lastAlert, exist := m.smsCache[deviceUUIDParam]
		if !exist || now.After(lastAlert.Add(10*time.Minute)) {
			m.smsCache[deviceUUIDParam] = now
			log.Infoln(fmt.Sprintf("High CO Density: %f Device: %s", *CODensity.Density, deviceUUIDParam))
			go func() {
				if err = m.smsProvider.SendMessage(
					"2000500666",
					"09397737262",
					fmt.Sprintf("High CO Density: %f\nDevice: %s", *CODensity.Density, deviceUUIDParam),
				); err != nil {
					log.Warnf("error in sending sms: %s", err.Error())
				}
			}()
		}
	}

	// Store it in server metrics
	log.Debugf("Recieved a new CO Density value: %f from org: %s device: %s", *CODensity.Density, orgUUIDParam, deviceUUIDParam)
	metrics.COGauge.WithLabelValues(orgUUID.String(), deviceUUID.String()).Set(*CODensity.Density)

	c.JSON(http.StatusOK, gin.H{
		"message": "CO Density pushed successfully",
	})
}

func (m *Metrics) RegisterRoutes(group *gin.RouterGroup) {
	group.PUT("/organizations/:orgUUID/devices/:deviceUUID/metric/co_density", m.PushCODensity)
}
