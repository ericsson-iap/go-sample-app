// COPYRIGHT Ericsson 2024

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package logging

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"eric-oss-hello-world-go-app/src/internal/configuration"

	"github.com/sirupsen/logrus"
)

var logger struct {
	conf    *configuration.Config
	tlsConf *tls.Config
	logrus  *logrus.Logger
	client  *http.Client
	wg      sync.WaitGroup
	level   logrus.Level
}

const (
	// PanicLevel Panic 0
	PanicLevel logrus.Level = iota
	// FatalLevel Fatal 1
	FatalLevel
	// ErrorLevel Error 2
	ErrorLevel
	// WarningLevel Warning 3
	WarningLevel
	// InfoLevel Info 4
	InfoLevel
	// DebugLevel Debug 5
	DebugLevel
)

type logControl struct {
	Severity  string `json:"severity"`
	Container string `json:"container"`
}

type logEntry struct {
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	Message   string `json:"message"`
	ServiceID string `json:"service_id"`
	Severity  string `json:"severity"`
}

// Init Initialize Logger
func Init() {
	logger.logrus = logrus.New()
	SetOutput(os.Stdout)
	SetLevel(InfoLevel)
	logger.conf = configuration.AppConfig
	logger.tlsConf = configuration.LogmTLSConfig()
	logger.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: logger.tlsConf,
		},
	}
	data, err := os.ReadFile(logger.conf.LogControlFile)
	if err != nil {
		logger.logrus.Error(logger.conf.LogControlFile)
		logger.logrus.Error(err)
		logger.logrus.Warn("Could not read from LogControlFile, setting level to INFO")
		return
	}

	var logControls []logControl

	if err := json.Unmarshal(data, &logControls); err != nil {
		logger.logrus.Error(err)
		logger.logrus.Warn("Could not parse LogControlFile, setting level to INFO")
		return
	}

	for _, item := range logControls {
		if item.Container == logger.conf.ContainerName {
			switch item.Severity {
			case "critical":
				SetLevel(FatalLevel)
			case "error":
				SetLevel(ErrorLevel)
			case "warning":
				SetLevel(WarningLevel)
			case "info":
				SetLevel(InfoLevel)
			case "debug":
				SetLevel(DebugLevel)
			}
			break
		}
	}
}

// SetLevel Set Log Level
func SetLevel(level logrus.Level) {
	logger.level = level
	logger.logrus.SetLevel(level)
}

// SetOutput Set Logger Output
func SetOutput(writer io.Writer) {
	logger.logrus.SetOutput(writer)
}

// Error Log at Error level
func Error(msg string) {
	if logger.level >= ErrorLevel {
		logger.logrus.Error(msg)
		dispatch(msg, ErrorLevel)
	}
}

// Warning Log at Warning level
func Warning(msg string) {
	if logger.level >= WarningLevel {
		logger.logrus.Warning(msg)
		dispatch(msg, WarningLevel)
	}
}

// Info Log at Info level
func Info(msg string) {
	if logger.level >= InfoLevel {
		logger.logrus.Info(msg)
		dispatch(msg, InfoLevel)
	}
}

// Debug Log at Debug level
func Debug(msg string) {
	if logger.level >= DebugLevel {
		logger.logrus.Debug(msg)
		dispatch(msg, DebugLevel)
	}
}

func dispatch(msg string, level logrus.Level) {
	if logger.tlsConf == nil {
		return
	}

	entry := &logEntry{time.Now().Format(time.RFC3339), "0.0.1", msg, "rapp-eric-oss-hello-world-go-app", level.String()}
	entryJSON, _ := json.Marshal(entry)
	logger.wg.Wait()
	logger.wg.Add(1)
	go func() {
		defer logger.wg.Done()
		request, err := http.NewRequestWithContext(context.Background(), "POST",
			"https://"+logger.conf.LogEndpoint, bytes.NewBuffer(entryJSON))
		if err != nil {
			logger.logrus.Error("Request failed for mTLS logging")
			return
		}

		request.Header.Set("Content-Type", "application/json")

		resp, err := logger.client.Do(request)
		if err != nil {
			logger.logrus.Error("Request failed for mTLS logging")
			return
		}
		defer resp.Body.Close() //nolint:errcheck //error has no impact
	}()
}
