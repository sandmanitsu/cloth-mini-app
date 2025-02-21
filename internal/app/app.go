package app

import (
	congig "cloth-mini-app/internal/config"
	"cloth-mini-app/internal/delivery/rest"
	sl "cloth-mini-app/internal/logger"
	repository "cloth-mini-app/internal/repository/item"
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
	e.Static("/admin/static", "public")

	rest.NewItemHandler(e, itemService)
	rest.NewAdminHandler(e)

	logger.Info("echo", sl.Err(e.Start(config.Host+":"+config.Port)))
}
