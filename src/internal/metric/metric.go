package metric

import (
	"eric-oss-hello-world-go-app/src/internal/configuration"

	"strings"

	"github.com/prometheus/client_golang/prometheus"
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
			Namespace: strings.Replace(configuration.AppConfig.ChosenName, "-", "_", -1),
			Name:      "requests_total",
			Help:      "Total number of API requests",
		})
}

func registerMetrics() {
	Registry.Register(RequestsTotal) //nolint:errcheck // handling invalid metrics descriptors is outside the app scope
}

// SetupMetrics sets up the metrics
func SetupMetrics() {
	createMetrics()
	registerMetrics()
}
