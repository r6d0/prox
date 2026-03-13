package main

import (
	"log/slog"
	"os"
	"os/signal"
	prox "prox/internal"
	"syscall"
)

func main() {
	var err error
	var config *prox.ProxConfig
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	if len(os.Args) > 1 {
		if _, checkErr := os.Stat(os.Args[1]); checkErr == nil {
			config, err = prox.NewJsonConfig(os.Args[1])
		}
	}

	if err == nil && config == nil {
		slog.Info("The configuration was not found. Prox will use the default configuration.")
		config = prox.NewDefaultConfig()
	}

	var server *prox.Prox
	if err == nil {
		if server, err = prox.NewProx(config); err == nil {
			go server.Start()
		}
	}

	if err == nil {
		slog.Info("Prox is running successfully.")

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop // Stop signal waiting.

		err = server.Stop()
	}

	if err != nil {
		panic(err.Error())
	} else {
		slog.Info("Proxy was successfully stopped.")
	}
}
