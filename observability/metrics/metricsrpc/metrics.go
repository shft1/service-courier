package metricsrpc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type RPCMetrics struct {
	Retry *prometheus.CounterVec
}

func NewRPCMetrics() *RPCMetrics {
	return &RPCMetrics{
		Retry: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "gateway_retries_total",
			Help: "Total number of gateway retries",
		}, []string{"name"}),
	}
}
