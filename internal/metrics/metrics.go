package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	TempGauge *prometheus.GaugeVec
	WindSpeedGauge
}

func NewMetrics(registry prometheus.Registerer) *Metrics {
	tempGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "weather",
			Name:      "temperature",
			Help:      "Gauge Vector indicating location temperature",
		},
		[]string{
			"location",
			"name",
		},
	)

	registry.MustRegister(tempGauge)

	return &Metrics{TempGauge: tempGauge}
}
