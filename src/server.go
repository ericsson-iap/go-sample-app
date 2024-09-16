package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eric-oss-hello-world-go-app/src/internal/configuration"
	log "eric-oss-hello-world-go-app/src/internal/logging"
	"eric-oss-hello-world-go-app/src/internal/metric"
	"eric-oss-hello-world-go-app/src/internal/request"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	config     = configuration.AppConfig
	server     *http.Server
	ExitSignal chan os.Signal
)

func init() {
	log.Init()
	log.Info("Go Hello World Sample App initializing...")
	ExitSignal = getExitSignal()
	metric.SetupMetrics()
}

func hello(resp http.ResponseWriter, req *http.Request) {

	metric.RequestsTotal.Inc()

	err := request.HandleLogin(config.IamClientID, config.IamClientSecret, config.IamBaseURL)
	if err != nil {
		log.Error("login failed: " + err.Error())
	}

	_, err = fmt.Fprintf(resp, "Hello World!!")
	if err != nil {
		log.Error("Error writing to response")
	}

	log.Info("Hello World!!")
}

func health(resp http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(resp, "Ok")
	if err != nil {
		log.Error("Error writing to response")
	}
	log.Debug("Health check: Ok")
}

// make one channel out of these termination signals so we can wait on one signal to exit the app
func getExitSignal() chan os.Signal {
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

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(metric.Registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/health", health)

	localPort := fmt.Sprintf(":%d", config.LocalPort)

	server = &http.Server{
		Addr:              localPort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
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

	log.Info("Server is ready to receive web requests")

	return server
}

func main() {
	go startWebService()
	<-ExitSignal //wait to receive exit signal
}
