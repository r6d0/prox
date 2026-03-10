package internal

import (
	"log"
	"os"
)

// The logging level.
type LogLevel uint8

const (
	ERROR LogLevel = iota
	INFO
	DEBUG
)

// The wrapper of the standard log.Logger.
type LoggerWrapper struct {
	Logger *log.Logger
	Level  LogLevel
}

// The function prints a message when the level is ERROR.
func (wrp *LoggerWrapper) Error(values ...any) {
	if wrp.Level <= ERROR {
		wrp.Logger.Println(values...)
	}
}

// The function prints a message when the level is INFO.
func (wrp *LoggerWrapper) Info(values ...any) {
	if wrp.Level <= INFO {
		wrp.Logger.Println(values...)
	}
}

// The function prints a message when the level is DEBUG.
func (wrp *LoggerWrapper) Debug(values ...any) {
	if wrp.Level <= DEBUG {
		wrp.Logger.Println(values...)
	}
}

// The function creates a new instance of LoggerWrapper.
func NewLoggerWrapper(config *ProxConfig) *LoggerWrapper {
	level := INFO
	if config != nil {
		level = config.Log.Level
	}
	return &LoggerWrapper{Logger: log.New(os.Stdout, "", log.LstdFlags), Level: level}
}
