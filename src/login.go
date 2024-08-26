// COPYRIGHT Ericsson 2024

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// Token retrieved here
type Token struct {
	AccessToken string `json:"accessToken"`
}

const loginPath = "/auth/realms/master/protocol/openid-connect/token"

// HandleLogin Creates an instance of the request body
func HandleLogin(clientID, clientSecret, baseURL string) error {
	loginURL := baseURL + path.Join(loginPath)
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

// CreateFormData Creates a formData Map that will be used with HandleFormRequest
func CreateFormData(clientID, clientSecret string) url.Values {
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)
	formData.Set("tenant_id", "master")
	return formData
}
