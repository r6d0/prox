package main

import (
	"errors"
	"os"
	"os/signal"
	prox "prox/internal"
	"syscall"
)

func main() {
	logger := prox.NewLoggerWrapper(nil)

	var err error
	var config *prox.ProxConfig
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	path := prox.DEFAULT_JSON_CONFIG
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		logger.Info("The configuration was not found. Prox will use the default configuration.")
		err = nil
		config = prox.NewDefaultConfig()
	} else if err == nil {
		config, err = prox.NewJsonConfig(path)
	}

	var server prox.Prox
	if err == nil {
		if server, err = prox.NewProx(config); err == nil {
			go server.Start()
		}
	}

	if err == nil {
		logger.Info("Prox is running successfully.")

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop // Stop signal waiting.

		err = server.Stop()
	}

	if err != nil {
		panic(err.Error())
	} else {
		logger.Info("Proxy was successfully stopped.")
	}
}
