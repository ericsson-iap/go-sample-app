package network

import (
	"net/http"
)

// GetIPInfo returns concatenated Ips from X-Forwarded-For header and request remoteAddr
func GetIPInfo(request *http.Request) string {
	// Get comma seperated IPs, more details: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For
	xForwardedForIps := request.Header.Get("X-Forwarded-For")
	remoteAddr := request.RemoteAddr

	return "X-Forwarded-For: '" + xForwardedForIps + "', RemoteAddr: '" + remoteAddr + "'"
}
