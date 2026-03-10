package internal

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

const DEFAULT_PROX_PORT = 9999
const DEFAULT_JSON_CONFIG = "./proxconfig.json"

var ErrUnsupportedProtocol = errors.New("protocol is unsupported")

// The Protocol of the proxy server.
type Protocol uint8

const (
	HTTP Protocol = iota
	HTTPS
)

// The configuration of the proxy server.
type ProxConfig struct {
	Port     uint16
	Protocol Protocol
	Request  HttpRequestProxConfig
	Tls      TlsProxConfig
	Log      LogProxConfig
}

// The configuration of HTTP requests.
type HttpRequestProxConfig struct {
	Timeout         time.Duration
	ForwardedHeader bool
}

// The TLS configuration of the proxy server.
type TlsProxConfig struct {
	CertFile string
	KeyFile  string
}

// The logger configuration.
type LogProxConfig struct {
	Level LogLevel
}

// The function reads the configuration from a JSON file.
func NewJsonConfig(jsonFile string) (*ProxConfig, error) {
	data, err := os.ReadFile(jsonFile)
	if err == nil {
		config := &ProxConfig{}
		return config, json.Unmarshal(data, config)
	}
	return nil, err
}

// The function returns the default configuration.
func NewDefaultConfig() *ProxConfig {
	return &ProxConfig{
		Port:     DEFAULT_PROX_PORT,
		Protocol: HTTP,
		Request: HttpRequestProxConfig{
			Timeout:         2 * time.Second,
			ForwardedHeader: false,
		},
	}
}
