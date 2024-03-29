package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Observable interface {
	ObserveSuccess(float64, ...string)
	ObserveAPIError(string)
}

type Metrics struct {
	TempGauge             *prometheus.GaugeVec
	ApiBadResponseCounter *prometheus.CounterVec
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

	apiBadResponseCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "weather",
			Name:      "badResponse",
			Help:      "Counter Vector indicating bad api response",
		},
		[]string{
			"statusCode",
		},
	)

	registry.MustRegister(tempGauge, apiBadResponseCounter)

	return &Metrics{
		TempGauge:             tempGauge,
		ApiBadResponseCounter: apiBadResponseCounter,
	}
}

func (m *Metrics) ObserveSuccess(temp float64, lvs ...string) {
	m.TempGauge.WithLabelValues(lvs...).Set(temp)
}

func (m *Metrics) ObserveAPIError(lv string) {
	m.ApiBadResponseCounter.WithLabelValues(lv).Inc()
}
