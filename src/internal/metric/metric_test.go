package metric_test

import (
	"eric-oss-hello-world-go-app/src/internal/metric"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMetrics(t *testing.T) {
	t.Parallel()

	assert.Nil(t, metric.RequestsTotal)
	assert.Nil(t, metric.RequestsFailedTotal)
	assert.Nil(t, metric.HelloWorldHTTPRequestsTotal)

	metric.SetupMetrics()

	assert.NotNil(t, metric.RequestsTotal,
		"RequestsTotal has not been initialized")
	assert.NotNil(t, metric.RequestsFailedTotal,
		"RequestsFailedTotal has not been initialized")
	assert.NotNil(t, metric.HelloWorldHTTPRequestsTotal,
		"HelloWorldHTTPRequestsTotal has not been initialized")
}

func TestRegisterMetrics(t *testing.T) {
	t.Parallel()

	metric.SetupMetrics()

	metrics, err := metric.Registry.Gather()
	assert.NoError(t, err)
	assert.Len(t, metrics, 2)
}
