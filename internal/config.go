package internal

import (
	"encoding/json"
	"log/slog"
	"os"
	"prox/internal/rule"
	"strings"
	"time"
)

const DEFAULT_PROX_PORT = 9999

// The target for logs data.
type LogTargetType uint8

const (
	CONSOLE LogTargetType = iota
	FILE
)

// The alias of time.Duration
type TimeDuration time.Duration

func (tmt *TimeDuration) UnmarshalJSON(data []byte) error {
	text := strings.ReplaceAll(string(data), "\"", "")
	duration, err := time.ParseDuration(text)
	if err == nil {
		*tmt = TimeDuration(duration)
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
	Rules      rule.RequestRulesConfig `json:"rules"`
	BufferSize int                     `json:"bufferSize"`
	Timeout    TimeDuration            `json:"timeout"`
	Forwarded  bool                    `json:"forwardedHeader"`
}

// The logger configuration.
type LogProxConfig struct {
	Level  slog.Level        `json:"level"`
	Target []LogTargetType   `json:"target"`
	File   LogFileProxConfig `json:"file"`
}

// The configuration of the log files.
type LogFileProxConfig struct {
	Dir        string       `json:"dir"`
	Size       int64        `json:"size"`
	TimeToLive TimeDuration `json:"timeToLive"`
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
			Timeout:   TimeDuration(2 * time.Second),
			Forwarded: false,
		},
		Log: LogProxConfig{Level: slog.LevelDebug},
	}
}
