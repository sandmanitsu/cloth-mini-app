package app

import (
	"cloth-mini-app/internal/congig"
	"cloth-mini-app/internal/delivery/rest"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/repository"
	"cloth-mini-app/internal/service/item"
	"cloth-mini-app/internal/storage/postgresql"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

// Running application
func Run(config *congig.Config, logger *slog.Logger) {
	logger.Info("starting app...")

	storage, err := postgresql.NewPostgreSQL(config.DB)
	if err != nil {
		logger.Error("failed to init postgresql storage", sl.Err(err))
		os.Exit(1)
	}

	// prepare repositories
	itemRepo := repository.NewItemRepository(logger, storage)

	// prepare services
	itemService := item.NewItemService(logger, itemRepo)

	e := echo.New()

	rest.NewItemHandler(e, itemService)

	logger.Info("echo", sl.Err(e.Start(config.Host+":"+config.Port)))
}
