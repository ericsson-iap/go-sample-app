package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"eric-oss-hello-world-go-app/src/internal/configuration"

	log "eric-oss-hello-world-go-app/src/internal/logging"

	"github.com/stretchr/testify/assert"
)

const logOutputFileName = "testlogfile"

func TestHelloAndHealthResponseAreValid(t *testing.T) {
	tests := []struct {
		name         string
		endpoint     string
		status       int
		response     string
		routeHandler func(resp http.ResponseWriter, req *http.Request)
	}{
		{
			name: "`/hello` Hello Endpoint handler", endpoint: "/hello", status: 200, response: "Hello World!!",
			routeHandler: hello,
		},
		{
			name: "`/health` Health Endpoint handler", endpoint: "/health", status: 200, response: "Ok",
			routeHandler: health,
		},
	}

	for _, testParameters := range tests {
		// arrange
		request := httptest.NewRequest(http.MethodGet, testParameters.endpoint, nil)
		response := httptest.NewRecorder()

		// act
		testParameters.routeHandler(response, request)

		// assert
		res := response.Result()
		defer res.Body.Close() //nolint:errcheck //error has no impact
		assert.Equal(t, testParameters.status, res.StatusCode,
			"Status code should be 200, but got : "+strconv.Itoa(res.StatusCode))

		data, _ := io.ReadAll(res.Body)
		assert.NotNil(t, data, "Data should not be nill")

		assert.Equal(t, testParameters.response, string(data),
			fmt.Sprintf("Should be returned `%v` but got : %v", testParameters.response, string(data)))
	}
}

func TestGetHelloAndHealthEndPointReturnValidResponse(t *testing.T) {
	tests := []struct {
		name         string
		endpoint     string
		status       int
		response     string
		routeHandler func(resp http.ResponseWriter, req *http.Request)
	}{
		{
			name: "`/hello` route testing", endpoint: "/hello", status: 200, response: "Hello World!!",
			routeHandler: hello,
		},
		{
			name: "`/health` route testing", endpoint: "/health", status: 200, response: "Ok",
			routeHandler: health,
		},
	}

	for _, testParameters := range tests {
		// arrange
		router := http.NewServeMux()
		router.HandleFunc(testParameters.endpoint, testParameters.routeHandler)
		svr := httptest.NewServer(router)
		defer svr.Close()

		// act
		res, err := http.Get(fmt.Sprintf("%s%s", svr.URL, testParameters.endpoint))

		// assert
		assert.Nil(t, err, fmt.Sprintf("Could not send GET request:  %v", err))
		defer res.Body.Close() //nolint:errcheck //error has no impact
		assert.Equal(t, res.StatusCode, http.StatusOK,
			fmt.Sprintf("Expected Status Ok; got %v", res.Status))

		data, _ := io.ReadAll(res.Body)
		assert.Equal(t, testParameters.response, string(data),
			fmt.Sprintf("Expected response is `%v`, but got %v", testParameters.response, string(data)))
	}
}

func TestExitSignalChannel(t *testing.T) {
	channel := getExitSignal()

	assert.NotNil(t, channel, "Channel should not be nil")
	assert.Equal(t, 1, cap(channel), "Capacity should be 1")
	assert.Equal(t, "chan", reflect.ValueOf(channel).Kind().String(), "Kind should be 'chan'")
}

func TestStartWebService(t *testing.T) {
	retries := 3
	var res *http.Response
	var err error

	file, err := os.OpenFile(logOutputFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	assert.Nil(t, err, fmt.Sprintf("error creating file: %v", err))
	defer file.Close() //nolint:errcheck //error has no impact
	wrt := io.MultiWriter(os.Stdout, file)
	log.SetOutput(wrt)
	log.SetLevel(log.DebugLevel)
	// act
	srv := startWebService()

	// assert
	assert.NotNil(t, srv, "Should not be nill")

	for retries > 0 {
		client := http.Client{}
		res, err = client.Get("http://localhost:8050/hello")
		if err != nil {
			log.Error("Request failed")
			retries--
		} else {
			defer res.Body.Close() //nolint:errcheck //error has no impact
			break
		}
	}

	assert.NotNil(t, res, "Response should not be nill")
	assert.Equal(t, res.StatusCode, http.StatusOK, fmt.Sprintf("Expected Status Ok; got %v", res.Status))

	content, _ := io.ReadAll(res.Body)
	assert.NotNil(t, content, "Content should not be nill")
	assert.Equal(t, "Hello World!!", string(content), "Should be returned `Hello World!!` but got : "+string(content))

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = srv.Shutdown(ctxShutDown)
	assert.Nil(t, err)
	_ = os.Remove(logOutputFileName)
}

func TestStartWebServiceWithHttps(t *testing.T) {
	retries := 3
	var res *http.Response
	var err error

	config = &configuration.Config{
		LocalPort:       8050,
		LocalProtocol:   "https",
		CertFile:        "certificate.pem",
		KeyFile:         "key.pem",
		LogControlFile:  "test.file",
		LogEndpoint:     "",
		AppKey:          "test.key",
		AppCert:         "test.cert",
		AppCertFilePath: "/etc/tls/log/",
	}

	file, err := os.OpenFile(logOutputFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	assert.Nil(t, err, fmt.Sprintf("error creating file: %v", err))
	defer file.Close() //nolint:errcheck //error has no impact
	wrt := io.MultiWriter(os.Stdout, file)
	log.SetOutput(wrt)
	log.SetLevel(log.DebugLevel)
	// act
	srv := startWebService()

	// assert
	assert.NotNil(t, srv, "Should not be nil")

	for retries > 0 {
		client := http.Client{}
		res, err = client.Get("https://localhost:8050/hello")
		if err != nil {
			log.Error("Request failed")
			retries--
		} else {
			defer res.Body.Close() //nolint:errcheck //error has no impact
			break
		}
	}

	assert.Nil(t, res, "Response should be nil")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = srv.Shutdown(ctxShutDown)
	assert.Nil(t, err)
	_ = os.Remove(logOutputFileName)
}
