// COPYRIGHT Ericsson 2024

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"eric-oss-hello-world-go-app/src/internal/configuration"
)

// HandleFormRequest for Client Credential Flow Login
func HandleFormRequest(endpoint string, formData url.Values, headers http.Header) ([]byte, error) {
	// Create a new TLS config with the server's CA cert
	tlsConfig := configuration.NewTLSConfig()

	// Create an HTTP client with the custom TLS config
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Create a new http.Request object
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost, endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Create http.Request object failed: %w", err)
	}

	// Set the headers on the request
	if headers != nil {
		req.Header = headers
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request to the specified endpoint
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request Failed with following error: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck //error has no impact

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Reading response body failed: %w", err)
	}

	// Check the response status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &httpError{
			statusCode: resp.StatusCode,
			statusText: resp.Status,
			body:       respBody,
		}
	}

	// If the response body is empty, return nil
	if len(respBody) == 0 {
		return nil, nil
	}

	// Return the response body
	return respBody, nil
}

type httpError struct {
	statusCode int
	statusText string
	body       []byte
}

func (e *httpError) Error() string {
	return e.statusText
}
