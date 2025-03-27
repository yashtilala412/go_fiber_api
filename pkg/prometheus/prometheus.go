package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const Namespace = "golang_api"

type PrometheusMetrics struct {
	RequestsMetrics *prometheus.CounterVec
}

var metrics *PrometheusMetrics = nil

func InitPrometheusMetrics() *PrometheusMetrics {
	if metrics == nil {
		metrics = &PrometheusMetrics{
			RequestsMetrics: promauto.NewCounterVec(prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "requests_total",
				Help:      "Total API requests",
			}, []string{"code"}),
		}
	}

	return metrics
}
