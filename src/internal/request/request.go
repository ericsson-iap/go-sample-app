package request

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"eric-oss-hello-world-go-app/src/internal/configuration"
)

// Token retrieved here
type Token struct {
	AccessToken string `json:"accessToken"`
}

const loginPath = "/auth/realms/master/protocol/openid-connect/token"

// HandleLogin Creates an instance of the request body
func HandleLogin(clientID, clientSecret, baseURL string) error {
	loginURL := baseURL + path.Join(loginPath)

	if len(clientID) == 0 || len(clientSecret) == 0 {
		return fmt.Errorf("Empty parameters provided for IamClientID or IamClientSecret")
	}
	formData := CreateFormData(clientID, clientSecret)

	respBody, err := HandleFormRequest(loginURL, formData, http.Header{})
	if err != nil {
		return err
	}
	var token Token
	if err := json.Unmarshal(respBody, &token); err != nil {
		return fmt.Errorf("JSON Unmarshal Failed with following error: %w", err)
	}

	return nil
}

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

// CreateFormData Creates a formData Map that will be used with HandleFormRequest
func CreateFormData(clientID, clientSecret string) url.Values {
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)
	formData.Set("tenant_id", "master")
	return formData
}

type httpError struct {
	statusCode int
	statusText string
	body       []byte
}

func (e *httpError) Error() string {
	return e.statusText
}
