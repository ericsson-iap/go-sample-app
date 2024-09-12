package logging

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"eric-oss-hello-world-go-app/src/internal/configuration"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const logOutputFileName = "testlogfile"

func createTestFile(t *testing.T, fileName, fileContent string) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	assert.Nil(t, err, fmt.Sprintf("error creating file: %v", err))
	n, err := file.WriteString(fileContent)
	assert.NotEqual(t, n, 0)
	assert.Nil(t, err)
	err = file.Close()
	assert.Nil(t, err)
	t.Cleanup(func() {
		e := os.Remove(fileName)
		assert.Nil(t, e, fmt.Sprintf("error deleting file: %v", e))
	})
}

func TestInitWithoutLogCtrl(t *testing.T) {
	Init()
	assert.NotNil(t, logger.logrus)
	assert.Equal(t, logger.level, InfoLevel)
	assert.Equal(t, logger.logrus.GetLevel(), InfoLevel)
}

func TestInitWithLogCtrl(t *testing.T) {
	createTestFile(t, "logcontrol.json", "[{\"severity\": \"debug\",\"container\": \"rapp-eric-oss-hello-world-go-app\"}]")
	t.Setenv("LOG_CTRL_FILE", "logcontrol.json")
	t.Setenv("CONTAINER_NAME", "rapp-eric-oss-hello-world-go-app")
	configuration.ReloadAppConfig()
	Init()
	assert.True(t, true)
	assert.NotNil(t, logger.logrus)
	assert.Equal(t, logger.level, DebugLevel)
	assert.Equal(t, logger.logrus.GetLevel(), DebugLevel)
}

func TestInitWithInfoLevel(t *testing.T) {
	createTestFile(t, "logcontrol.json", "[{\"severity\": \"info\",\"container\": \"rapp-eric-oss-hello-world-go-app\"}]")
	t.Setenv("LOG_CTRL_FILE", "logcontrol.json")
	t.Setenv("CONTAINER_NAME", "rapp-eric-oss-hello-world-go-app")
	configuration.ReloadAppConfig()
	Init()
	assert.True(t, true)
	assert.NotNil(t, logger.logrus)
	assert.Equal(t, logger.level, InfoLevel)
	assert.Equal(t, logger.logrus.GetLevel(), InfoLevel)
}

func TestInitWithWarningLevel(t *testing.T) {
	createTestFile(t, "logcontrol.json", "[{\"severity\": \"warning\",\"container\": \"rapp-eric-oss-hello-world-go-app\"}]")
	t.Setenv("LOG_CTRL_FILE", "logcontrol.json")
	t.Setenv("CONTAINER_NAME", "rapp-eric-oss-hello-world-go-app")
	configuration.ReloadAppConfig()
	Init()
	assert.True(t, true)
	assert.NotNil(t, logger.logrus)
	assert.Equal(t, logger.level, WarningLevel)
	assert.Equal(t, logger.logrus.GetLevel(), WarningLevel)
}

func TestInitWithErrorLevel(t *testing.T) {
	createTestFile(t, "logcontrol.json", "[{\"severity\": \"error\",\"container\": \"rapp-eric-oss-hello-world-go-app\"}]")
	t.Setenv("LOG_CTRL_FILE", "logcontrol.json")
	t.Setenv("CONTAINER_NAME", "rapp-eric-oss-hello-world-go-app")
	configuration.ReloadAppConfig()
	Init()
	assert.True(t, true)
	assert.NotNil(t, logger.logrus)
	assert.Equal(t, logger.level, ErrorLevel)
	assert.Equal(t, logger.logrus.GetLevel(), ErrorLevel)
}

func TestInitWithFatalLevel(t *testing.T) {
	createTestFile(t, "logcontrol.json", "[{\"severity\": \"critical\",\"container\": \"rapp-eric-oss-hello-world-go-app\"}]")
	t.Setenv("LOG_CTRL_FILE", "logcontrol.json")
	t.Setenv("CONTAINER_NAME", "rapp-eric-oss-hello-world-go-app")
	configuration.ReloadAppConfig()
	Init()
	assert.True(t, true)
	assert.NotNil(t, logger.logrus)
	assert.Equal(t, logger.level, FatalLevel)
	assert.Equal(t, logger.logrus.GetLevel(), FatalLevel)
}

func TestInitWithInvalidLogCtrl(t *testing.T) {
	createTestFile(t, "logcontrol.json", "[][]")
	t.Setenv("LOG_CTRL_FILE", "logcontrol.json")
	t.Setenv("CONTAINER_NAME", "rapp-eric-oss-hello-world-go-app")
	configuration.ReloadAppConfig()
	Init()
	assert.True(t, true)
	assert.NotNil(t, logger.logrus)
	assert.Equal(t, logger.level, InfoLevel)
	assert.Equal(t, logger.logrus.GetLevel(), InfoLevel)
}

func TestSetLevel(t *testing.T) {
	Init()
	assert.NotNil(t, logger.logrus)
	assert.Equal(t, logger.level, InfoLevel)
	assert.Equal(t, logger.logrus.GetLevel(), InfoLevel)

	SetLevel(DebugLevel)
	assert.Equal(t, logger.level, DebugLevel)
	assert.Equal(t, logger.logrus.GetLevel(), DebugLevel)
}

func TestLogOnLowestLevel(t *testing.T) {
	Init()
	assert.NotNil(t, logger.logrus)
	assert.Equal(t, logger.level, InfoLevel)
	assert.Equal(t, logger.logrus.GetLevel(), InfoLevel)

	file, err := os.OpenFile(logOutputFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	assert.Nil(t, err, fmt.Sprintf("error creating file: %v", err))

	wrt := io.MultiWriter(os.Stdout, file)
	SetLevel(DebugLevel)
	SetOutput(wrt)

	Error("error test")
	Warning("warning test")
	Info("info test")
	Debug("debug test")

	data, err := os.ReadFile(logOutputFileName)
	assert.Nil(t, err, fmt.Sprintf("error opening file: %v", err))

	assert.True(t, strings.Contains(string(data), "error test"), "This message should be logged")
	assert.True(t, strings.Contains(string(data), "warning test"), "This message should be logged")
	assert.True(t, strings.Contains(string(data), "info test"), "This message should be logged")
	assert.True(t, strings.Contains(string(data), "debug test"), "This message should be logged")
	t.Cleanup(func() {
		err := file.Close()
		assert.Nil(t, err)
		e := os.Remove(logOutputFileName)
		assert.Nil(t, e, fmt.Sprintf("error deleting test log file: %v", e))
	})
}

func TestLogOnHighestLevel(t *testing.T) {
	Init()
	assert.NotNil(t, logger.logrus)
	assert.Equal(t, logger.level, InfoLevel)
	assert.Equal(t, logger.logrus.GetLevel(), InfoLevel)

	file, err := os.OpenFile(logOutputFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	assert.Nil(t, err, fmt.Sprintf("error creating file: %v", err))

	wrt := io.MultiWriter(os.Stdout, file)
	SetLevel(FatalLevel)
	SetOutput(wrt)

	Error("error test")
	Warning("warning test")
	Info("info test")
	Debug("debug test")

	data, err := os.ReadFile(logOutputFileName)
	assert.Nil(t, err, fmt.Sprintf("error opening file: %v", err))

	assert.False(t, strings.Contains(string(data), "error test"), "This message should not be logged")
	assert.False(t, strings.Contains(string(data), "warning test"), "This message should not be logged")
	assert.False(t, strings.Contains(string(data), "info test"), "This message should not be logged")
	assert.False(t, strings.Contains(string(data), "debug test"), "This message should not be logged")
	t.Cleanup(func() {
		err = file.Close()
		assert.Nil(t, err)
		e := os.Remove(logOutputFileName)
		assert.Nil(t, e, fmt.Sprintf("error deleting test log file: %v", e))
	})
}

func TestDispatch(t *testing.T) {
	err := generateCACert()
	assert.Nil(t, err, fmt.Sprintf("error generating cacert: %v", err))
	err = generateKeyCertPair()
	assert.Nil(t, err, fmt.Sprintf("error generating certs: %v", err))
	t.Setenv("CA_CERT_FILE_NAME", "cacert.crt")
	t.Setenv("APP_CERT", "cert.pem")
	t.Setenv("APP_KEY", "key.pem")
	configuration.ReloadAppConfig()
	Init()
	Info("test")
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
