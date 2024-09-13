package network_test

import (
	"eric-oss-hello-world-go-app/src/internal/network"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIpAddressWithOnlyRemoteAddr(t *testing.T) {
	t.Parallel()

	// arrange
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	request.RemoteAddr = "192.0.0.1:8080"

	// act
	ipAddress := network.GetIPInfo(request)

	// assert
	assert.NotNil(t, ipAddress, "IP address should not be nil")
	assert.Contains(t, ipAddress, "RemoteAddr: '192.0.0.1:8080'")
}

func TestIpAddressWithOnlyXForwardedForHeaderSet(t *testing.T) {
	t.Parallel()

	// arrange
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	request.RemoteAddr = ""
	request.Header.Set("X-Forwarded-For", "192.0.0.1,192.0.0.2")

	// act
	ipAddress := network.GetIPInfo(request)

	// assert
	assert.NotNil(t, ipAddress, "IP address should not be nil")
	assert.Contains(t, ipAddress, "X-Forwarded-For: '192.0.0.1,192.0.0.2'")
}

func TestIpAddressWithBothRemoteAddrAndXForwardedForHeader(t *testing.T) {
	t.Parallel()

	// arrange
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	request.RemoteAddr = "192.79.12.10:9090"
	request.Header.Set("X-Forwarded-For", "192.0.0.1,192.0.0.2")

	// act
	ipAddress := network.GetIPInfo(request)

	// assert
	assert.Equal(t, ipAddress, "X-Forwarded-For: '192.0.0.1,192.0.0.2', RemoteAddr: '192.79.12.10:9090'")
}

func TestIpAddressWithoutXForwardedForHeaderAndRemoteAddr(t *testing.T) {
	t.Parallel()

	// arrange
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	request.Header.Set("X-Forwarded-For", "")
	request.RemoteAddr = ""

	// act
	ipAddress := network.GetIPInfo(request)

	// assert
	assert.Equal(t, "X-Forwarded-For: '', RemoteAddr: ''", ipAddress)
}
