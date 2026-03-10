package internal

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

const DEFAULT_PROX_PORT = 9999
const DEFAULT_JSON_CONFIG = "./proxconfig.json"

var ErrUnsupportedProtocol = errors.New("protocol is unsupported")

// The Protocol of the proxy server.
type Protocol uint8

func (prt Protocol) String() string {
	switch prt {
	case HTTP:
		return "HTTP"
	case HTTPS:
		return "HTTPS"
	}
	panic("unexpected Protocol value")
}

func (prt *Protocol) UnmarshalJSON(data []byte) error {
	var err error

	text := strings.ReplaceAll(string(data), "\"", "")
	if strings.EqualFold(text, HTTP.String()) {
		*prt = HTTP
	} else if strings.EqualFold(text, HTTPS.String()) {
		*prt = HTTPS
	} else {
		err = errors.New("unexpected Protocol value: " + text)
	}
	return err
}

const (
	HTTP Protocol = iota
	HTTPS
)

// The alias of time.Duration
type Timeout time.Duration

func (tmt *Timeout) UnmarshalJSON(data []byte) error {
	text := strings.ReplaceAll(string(data), "\"", "")
	duration, err := time.ParseDuration(text)
	if err == nil {
		*tmt = Timeout(duration)
	}
	return err
}

// The configuration of the proxy server.
type ProxConfig struct {
	Port     uint16                `json:"port"`
	Protocol Protocol              `json:"protocol"`
	Request  HttpRequestProxConfig `json:"request"`
	Tls      TlsProxConfig         `json:"tls"`
	Log      LogProxConfig         `json:"log"`
}

// The configuration of HTTP requests.
type HttpRequestProxConfig struct {
	Timeout   Timeout `json:"timeout"`
	Forwarded bool    `json:"forwardedHeader"`
}

// The TLS configuration of the proxy server.
type TlsProxConfig struct {
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

// The logger configuration.
type LogProxConfig struct {
	Level LogLevel `json:"level"`
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
			Timeout:   Timeout(2 * time.Second),
			Forwarded: false,
		},
	}
}
