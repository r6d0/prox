package internal

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

const DEFAULT_PROX_PORT = 9999

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
	Port    uint16                `json:"port"`
	Request HttpRequestProxConfig `json:"request"`
	Log     LogProxConfig         `json:"log"`
}

// The configuration of HTTP requests.
type HttpRequestProxConfig struct {
	Timeout   Timeout `json:"timeout"`
	Forwarded bool    `json:"forwardedHeader"`
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
		Port: DEFAULT_PROX_PORT,
		Request: HttpRequestProxConfig{
			Timeout:   Timeout(2 * time.Second),
			Forwarded: false,
		},
	}
}
