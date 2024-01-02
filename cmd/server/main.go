package main

import (
	"flag"

	"go.uber.org/zap"
)

func main() {
	var dev bool
	flag.BoolVar(&dev, "dev", false, "run in development environment")

	var configPath string
	flag.StringVar(&configPath, "config", "", "path to the configuration file")
	flag.Parse()

	if configPath == "" {
		panic("-config flag is not set")
	}

	var logger *zap.Logger
	var err error
	if dev {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	logger.Info("🚀 starting the application")
}
