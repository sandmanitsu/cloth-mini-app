package main

import (
	"cloth-mini-app/internal/app"
	congig "cloth-mini-app/internal/config"
	"cloth-mini-app/internal/logger"
	"log"
	"log/slog"
)

func main() {
	log.Println("config initializing...")
	config := congig.MustLoad()

	log.Println("logger initializing...")
	logger := logger.NewLogger(config.Env)
	logger.Info("logger started!", slog.String("env", config.Env))

	app.Run(config, logger)
}
