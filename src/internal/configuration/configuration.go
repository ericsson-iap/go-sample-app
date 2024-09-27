// Package configuration provides a pattern for the application to retrieve OS environment variables
package configuration

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path"
	"strconv"
	"strings"
)

// Config is a struct that contains all fields currently read from OS environment variables
type Config struct {
	LocalPort       int
	LocalProtocol   string
	CertFile        string
	KeyFile         string
	ContainerName   string
	IamClientID     string
	IamClientSecret string
	IamBaseURL      string
	CaCertFileName  string
	CaCertFilePath  string
	LogControlFile  string
	LogEndpoint     string
	AppKey          string
	AppCert         string
	AppCertFilePath string
}

const localPort = 8050

// AppConfig contains a list of values read from OS environment variables
var AppConfig = configFromEnvVars()

// ReloadAppConfig can be used to force a re-read of OS environment variables
func ReloadAppConfig() {
	AppConfig = configFromEnvVars()
}

func configFromEnvVars() *Config {
	return &Config{
		LocalPort:       getOsEnvInt("LOCAL_PORT", localPort),
		LocalProtocol:   getOsEnvString("LOCAL_PROTOCOL", "http"),
		CertFile:        getOsEnvString("CERT_FILE", "certificate.pem"),
		KeyFile:         getOsEnvString("KEY_FILE", "key.pem"),
		ContainerName:   getOsEnvString("CONTAINER_NAME", ""),
		IamClientID:     getOsEnvString("IAM_CLIENT_ID", ""),
		IamClientSecret: getOsEnvString("IAM_CLIENT_SECRET", ""),
		IamBaseURL:      getOsEnvString("IAM_BASE_URL", ""),
		CaCertFileName:  getOsEnvString("CA_CERT_FILE_NAME", ""),
		CaCertFilePath:  getOsEnvString("CA_CERT_FILE_PATH", ""),
		LogControlFile:  getOsEnvString("LOG_CTRL_FILE", ""),
		LogEndpoint:     getOsEnvString("LOG_ENDPOINT", ""),
		AppKey:          getOsEnvString("APP_KEY", ""),
		AppCert:         getOsEnvString("APP_CERT", ""),
		AppCertFilePath: getOsEnvString("APP_CERT_FILE_PATH", ""),
	}
}

func getOsEnvInt(envName string, defaultValue int) int {
	envValue := strings.TrimSpace(os.Getenv(envName))
	result, err := strconv.Atoi(envValue)
	if err != nil {
		result = defaultValue
	}

	return result
}

func getOsEnvString(envName, defaultValue string) string {
	result := strings.TrimSpace(os.Getenv(envName))

	if result == "" {
		result = defaultValue
	}

	return result
}

// NewTLSConfig Create new TLS configuration
func NewTLSConfig() *tls.Config {
	// Load the root CA certificate
	platformCaCert, err := os.ReadFile(getCertPath())
	if err != nil {
		return nil
	}

	// Create a new CertPool and add the root CA certificate to it
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(platformCaCert)

	// Create a TLS config with the client certificates and the server's CA cert
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{},
		RootCAs:            caCertPool,
		MinVersion:         tls.VersionTLS13,
	}

	return tlsConfig
}

// LogmTLSConfig Create new mTLS configuration for logging
func LogmTLSConfig() *tls.Config {
	// Load the root CA certificate
	platformCaCert, err := os.ReadFile(getCertPath())
	if err != nil {
		return nil
	}

	// Create a new CertPool and add the root CA certificate to it
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(platformCaCert)

	certFilePath := path.Join(AppConfig.AppCertFilePath, AppConfig.AppCert)
	keyFilePath := path.Join(AppConfig.AppCertFilePath, AppConfig.AppKey)

	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		return nil
	}
	// Create a TLS config with the client certificates and the server's CA cert
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		MinVersion:         tls.VersionTLS13,
	}

	return tlsConfig
}

// combines CaMountPath and CaCertFileName as a full path
func getCertPath() (certFilePath string) {
	config := AppConfig
	certFilePath = path.Join(config.CaCertFilePath, config.CaCertFileName)
	return certFilePath
}
