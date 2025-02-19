package app

import (
	"cloth-mini-app/internal/congig"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/postgresql"
	"log/slog"
)

// Running application
func Run(config *congig.Config, logger *slog.Logger) {
	logger.Info("starting app...")

	storage, err := postgresql.NewPostgreSQL(config.DB)
	if err != nil {
		logger.Error("failed to init postgresql storage", sl.Err(err))
	}

	logger.Info("storage connected")
	_ = storage
}
