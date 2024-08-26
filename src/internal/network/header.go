// COPYRIGHT Ericsson 2023

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

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
