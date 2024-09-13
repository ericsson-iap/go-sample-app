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
	// RequestsFailedTotal total number of API request failures
	RequestsFailedTotal prometheus.Counter
	// HelloWorldHTTPRequestsTotal total number of HTTP responses by status codes
	HelloWorldHTTPRequestsTotal *prometheus.CounterVec
)

func createMetrics() {
	RequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: servicePrefix,
			Name:      "requests_total",
			Help:      "Total number of API requests",
		})
	RequestsFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: servicePrefix,
			Name:      "requests_failed_total",
			Help:      "Total number of API requests failures",
		})
	HelloWorldHTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: servicePrefix,
			Name:      "hello_world_http_requests_total",
			Help:      "Total number of HTTP responses by status codes",
		},
		[]string{"code"})
}

func registerMetrics() {
	Registry.Register(RequestsTotal)               //nolint:errcheck // handling invalid metrics descriptors is outside the app scope
	Registry.Register(RequestsFailedTotal)         //nolint:errcheck // handling invalid metrics descriptors is outside the app scope
	Registry.Register(HelloWorldHTTPRequestsTotal) //nolint:errcheck // handling invalid metrics descriptors is outside the app scope
}

// SetupMetrics sets up the metrics
func SetupMetrics() {
	createMetrics()
	registerMetrics()
}
