// COPYRIGHT Ericsson 2024

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package configuration

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	key                string = "OS_ENV"
	defaultValueInt    int    = 1
	defaultValueString string = "string"
)

func TestGetConfig(t *testing.T) {
	t.Parallel()

	testConfig := AppConfig

	assert.NotNil(t, testConfig,
		"Instance should not be nil")

	assert.Equal(t, 8050, testConfig.LocalPort,
		"Port should be 8050, but got : "+strconv.Itoa(testConfig.LocalPort))

	assert.Equal(t, "httpDeliberateFail", testConfig.LocalProtocol,
		"Protocol should be http, but got : "+testConfig.LocalProtocol)

	assert.Equal(t, "certificate.pem", testConfig.CertFile,
		"Certificate file name should be `certificate.pem`, but got : "+testConfig.CertFile)

	assert.Equal(t, "key.pem", testConfig.KeyFile,
		"Keyfile file name should be `key.pem`, but got : "+testConfig.KeyFile)

	assert.Equal(t, "", testConfig.LogControlFile,
		"LogControlFile should be an empty string, but got : "+testConfig.LogControlFile)
}

func TestReloadAppConfig(t *testing.T) {
	t.Parallel()
	config1 := AppConfig
	ReloadAppConfig()
	config2 := AppConfig
	assert.NotSame(t, config1, config2,
		"AppConfig should be different pointer after ReloadAppConfig()")
}

func TestGetOsEnvIntSet(t *testing.T) {
	t.Setenv(key, "123")

	result := getOsEnvInt(key, defaultValueInt)
	assert.Equal(t, 123, result)
}

func TestGetOsEnvIntUnset(t *testing.T) {
	t.Parallel()

	result := getOsEnvInt(key, defaultValueInt)
	assert.Equal(t, defaultValueInt, result)
}

func TestGetOsEnvIntSetBadInt(t *testing.T) {
	t.Setenv(key, "abc")

	result := getOsEnvInt(key, defaultValueInt)
	assert.Equal(t, defaultValueInt, result)
}

func TestGetOsEnvStringSet(t *testing.T) {
	t.Setenv(key, "someValue")

	result := getOsEnvString(key, defaultValueString)
	assert.Equal(t, "someValue", result)
}

func TestGetOsEnvStringUnset(t *testing.T) {
	t.Parallel()

	result := getOsEnvString(key, defaultValueString)
	assert.Equal(t, defaultValueString, result)
}

func TestInvalidTLSConfigs(t *testing.T) {
	tlsConfig := NewTLSConfig()
	assert.Nil(t, tlsConfig)

	logMtlsConfig := LogmTLSConfig()
	assert.Nil(t, logMtlsConfig)
}

func TestTLSConfigs(t *testing.T) {
	err := generateCACert()
	assert.Nil(t, err, fmt.Sprintf("error generating cacert: %v", err))
	err = generateKeyCertPair()
	assert.Nil(t, err, fmt.Sprintf("error generating certs: %v", err))
	t.Setenv("CA_CERT_FILE_NAME", "cacert.crt")
	t.Setenv("APP_CERT", "cert.pem")
	t.Setenv("APP_KEY", "key.pem")
	ReloadAppConfig()

	tlsConfig := NewTLSConfig()
	assert.NotNil(t, tlsConfig)

	logMtlsConfig := LogmTLSConfig()
	assert.NotNil(t, logMtlsConfig)

	t.Cleanup(func() {
		e := os.Remove("cacert.crt")
		assert.Nil(t, e, fmt.Sprintf("error deleting cacert.crt: %v", e))
		e = os.Remove("cert.pem")
		assert.Nil(t, e, fmt.Sprintf("error deleting cert.pem: %v", e))
		e = os.Remove("key.pem")
		assert.Nil(t, e, fmt.Sprintf("error deleting key.pem: %v", e))
	})
}

func generateCACert() error {
	file, err := os.OpenFile("cacert.crt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return err
	}
	err = file.Close()

	return err
}

func generateKeyCertPair() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 64)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	keyUsage := x509.KeyUsageDigitalSignature
	keyUsage |= x509.KeyUsageKeyEncipherment

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"test"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(10000),

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	certOut, err := os.Create("cert.pem")
	if err != nil {
		return err
	}

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}

	if err := certOut.Close(); err != nil {
		return err
	}

	keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return err
	}

	err = keyOut.Close()

	return err
}
