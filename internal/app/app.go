package app

import (
	congig "cloth-mini-app/internal/config"
	"cloth-mini-app/internal/delivery/rest"
	sl "cloth-mini-app/internal/logger"
	brandRepo "cloth-mini-app/internal/repository/brand"
	categoryRepo "cloth-mini-app/internal/repository/category"
	imageRepo "cloth-mini-app/internal/repository/image"
	itemRepo "cloth-mini-app/internal/repository/item"
	"cloth-mini-app/internal/service/brand"
	"cloth-mini-app/internal/service/category"
	"cloth-mini-app/internal/service/image"
	"cloth-mini-app/internal/service/item"
	"cloth-mini-app/internal/storage/minio"
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

	minioClient, err := minio.NewMinioClient(config.Minio)
	if err != nil {
		logger.Error("failed to init minio storage", sl.Err(err))
		os.Exit(1)
	}
	_ = minioClient

	// prepare repositories
	itemRepo := itemRepo.NewItemRepository(logger, storage)
	categoryRepo := categoryRepo.NewCategoryRepository(logger, storage)
	brandRepo := brandRepo.NewBrandRepository(logger, storage)
	imageRepo := imageRepo.NewImageRepository(logger, storage)

	// prepare services
	itemService := item.NewItemService(logger, itemRepo, imageRepo)
	categoryService := category.NewCategoryService(logger, categoryRepo)
	brandService := brand.NewBrandService(logger, brandRepo)
	imageService := image.NewImageService(logger, minioClient, imageRepo)

	e := echo.New()
	e.Static("/admin/static", "public")

	// prepare handlers
	rest.NewItemHandler(e, itemService)
	rest.NewAdminHandler(e)
	rest.NewCategoryHandler(e, categoryService)
	rest.NewBrandHandler(e, brandService)
	rest.NewImageHandler(e, imageService)

	logger.Info("echo", sl.Err(e.Start(config.Host+":"+config.Port)))
}
