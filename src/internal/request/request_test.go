package request_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"eric-oss-hello-world-go-app/src/internal/request"

	"github.com/stretchr/testify/assert"
)

var formData = request.CreateFormData("ClientID", "ClientSecret")

func TestHandleLoginWithInvalidURL(t *testing.T) {
	// when invalid baseURL is passed, we are expecting the request to fail
	t.Parallel()
	err := request.HandleLogin("testID", "testSecret", "")
	assert.Contains(t, err.Error(), "Request Failed with following error: ")
}

func TestHandleLoginWithEmptyFormDataParameters(t *testing.T) {
	// when empty ClientID or ClientSecret is passed, we are expecting the request to fail
	t.Parallel()
	err := request.HandleLogin("", "", "")
	assert.Contains(t, err.Error(), "Empty parameters provided for IamClientID or IamClientSecret")
}

func TestHandleLoginWithJSON(t *testing.T) {
	// when server returns token, we are expecting no error to be returned
	t.Parallel()
	testResponse := request.Token{AccessToken: "testToken"}
	jsonResp, _ := json.Marshal(testResponse)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(jsonResp) //nolint:errcheck //mock server, no error handling required
	}))
	defer server.Close()

	err := request.HandleLogin("testID", "testSecret", server.URL)
	assert.Nil(t, err, "HandleLogin should return nil")
}

func TestHandleLoginWithoutJSON(t *testing.T) {
	// when server does not return JSON, we are expecting JSON Unmarshal error
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK`)) //nolint:errcheck //mock server, no error handling required
	}))
	defer server.Close()

	err := request.HandleLogin("testID", "testSecret", server.URL)
	assert.Contains(t, err.Error(), "JSON Unmarshal Failed with following error: ")
}

func TestHandleFormRequestWithResponse(t *testing.T) {
	// same response should be returned if server is sending response
	t.Parallel()
	testResponse := []byte(`test`)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(testResponse) //nolint:errcheck //mock server, no error handling required
	}))
	defer server.Close()
	resp, _ := request.HandleFormRequest(server.URL, formData, http.Header{})
	assert.Equal(t, resp, testResponse)
}

func TestHandleFormRequestWithNoResponse(t *testing.T) {
	// if no response returned, the len of byte returned is 0
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	}))
	defer server.Close()
	resp, _ := request.HandleFormRequest(server.URL, formData, http.Header{})
	assert.Equal(t, len(resp), 0)
}

func TestHandleFormRequestWithIncorrectStatusCode(t *testing.T) {
	// when receiving 403 response, error message should contain 403
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	_, err := request.HandleFormRequest(server.URL, formData, http.Header{})
	assert.Contains(t, err.Error(), "403")
}

func TestGetFormDataReturnsCorrectValues(t *testing.T) {
	// when method is called, we are expecting correct values to be returned
	t.Parallel()
	testID := "testID"
	testSecret := "testSecret"
	testFormData := request.CreateFormData(testID, testSecret)

	assert.NotNil(t, testFormData, "formData should not be nil")
	assert.Equal(t, testFormData.Get("client_id"), testID)
	assert.Equal(t, testFormData.Get("client_secret"), testSecret)
}
