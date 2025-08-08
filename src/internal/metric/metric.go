package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	servicePrefix = "hello_world"
)

var (
	// Registry prometheus registry
	Registry = prometheus.NewRegistry()
	// RequestsTotal total number of API requests
	RequestsTotal prometheus.Counter
)

func createMetrics() {
	RequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: servicePrefix,
			Name:      "requests_total",
			Help:      "Total number of API requests",
		})
}

func registerMetrics() {
	Registry.Register(RequestsTotal)               //nolint:errcheck // handling invalid metrics descriptors is outside the app scope
}

// SetupMetrics sets up the metrics
func SetupMetrics() {
	createMetrics()
	registerMetrics()
}
