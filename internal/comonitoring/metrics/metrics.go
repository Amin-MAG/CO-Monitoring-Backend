package metrics

import "github.com/prometheus/client_golang/prometheus"

func RegisterAll() error {
	appMetrics := []prometheus.Collector{
		COGauge,
	}

	// Register all metrics
	for _, m := range appMetrics {
		err := prometheus.Register(m)
		if err != nil {
			return err
		}
	}

	return nil
}
