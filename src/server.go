// COPYRIGHT Ericsson 2023

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package main

import (
	"context"
	"eric-oss-hello-world-go-app/src/internal/network"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"eric-oss-hello-world-go-app/src/internal/configuration"
	log "eric-oss-hello-world-go-app/src/internal/logging"
	"eric-oss-hello-world-go-app/src/internal/metric"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	config = configuration.AppConfig
	server *http.Server
	// ExitSignal The ExitSignal
	ExitSignal chan os.Signal
)

func handleAPICall(resp http.ResponseWriter, req *http.Request) {
	log.Debug("Entering api handler...")

	log.Info("Request IP: " + network.GetIPInfo(req))

	metric.RequestsTotal.Inc()

	err := HandleLogin(config.IamClientID, config.IamClientSecret, config.IamBaseURL)
	if err != nil {
		log.Error("Login Failed. " + err.Error())
	} else {
		log.Debug("Login Success.")
	}

	_, err = fmt.Fprintf(resp, "Hello World!!")
	if err != nil {
		log.Error("Error writing to response")
	}

	log.Debug("Leaving api handler...")
}

func checkServerHealth(resp http.ResponseWriter, req *http.Request) {
	// add some health checks here if required
	fmt.Fprintf(resp, "Ok")

}

func getExitSignalsChannel() chan os.Signal {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)
	return channel
}

func startWebService() *http.Server {
	ctx, servercancel := context.WithCancel(context.Background())

	initLogger()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(metric.Registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/hello", handleAPICall)
	mux.HandleFunc("/health", checkServerHealth)

	localPort := fmt.Sprintf(":%d", config.LocalPort)

	server = &http.Server{
		Addr:    localPort,
		Handler: mux,
	}

	go func() {
		if config.LocalProtocol == "https" {
			err := server.ListenAndServeTLS(config.CertFile, config.KeyFile)
			if err != nil {
				log.Error(err.Error())
			}
		} else {
			err := server.ListenAndServe()
			if err != nil {
				log.Error(err.Error())
			}
		}
		defer ctx.Done()
		defer servercancel()
	}()
	return server
}

func init() {
	log.Info("Hello World Sample App")
	ExitSignal = getExitSignalsChannel()

	metric.SetupMetrics()
}

// initLogger function to configure logger from
// log control file on server StartUp
func initLogger() {
	log.Init()
}

func main() {
	go startWebService()
	log.Info("Server is ready to receive web requests")

	<-ExitSignal
	log.Info("Terminating Hello World")
}
