package rest

import (
	"comonitoring/config"
	"comonitoring/internal/comonitoring/api/rest/health"
	"comonitoring/internal/comonitoring/api/rest/metric"
	"comonitoring/internal/comonitoring/api/rest/server"
	absms "comonitoring/internal/sms"
	"net/http"
	"time"
)

// NewApiServer creates the modules and the API server
func NewApiServer(c config.Config, smsProvider absms.SMS) (*http.Server, error) {
	smsCache := make(map[string]time.Time)

	// Create the health module
	healthMod, err := health.NewModule()
	if err != nil {
		return nil, err
	}

	// Create the metrics module
	metricMod, err := metric.NewModule(smsProvider, smsCache)
	if err != nil {
		return nil, err
	}

	return server.NewServer(healthMod, metricMod, c)
}
