package metrics

import "github.com/prometheus/client_golang/prometheus"

var COGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "systems_co_density",
	Help: "The density of the CO measured by system.",
}, []string{"organization_uuid", "device_uuid"})
