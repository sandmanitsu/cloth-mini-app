package app

import (
	"cloth-mini-app/internal/background"
	congig "cloth-mini-app/internal/config"
	"cloth-mini-app/internal/delivery/rest"
	"cloth-mini-app/internal/facade"
	"cloth-mini-app/internal/kafka"
	sl "cloth-mini-app/internal/logger"
	brandRepo "cloth-mini-app/internal/repository/brand"
	categoryRepo "cloth-mini-app/internal/repository/category"
	imageRepo "cloth-mini-app/internal/repository/image"
	itemRepo "cloth-mini-app/internal/repository/item"
	itemImageRepo "cloth-mini-app/internal/repository/item_image"
	lockRepo "cloth-mini-app/internal/repository/lock"
	outboxRepo "cloth-mini-app/internal/repository/outbox"
	"cloth-mini-app/internal/service/brand"
	"cloth-mini-app/internal/service/category"
	"cloth-mini-app/internal/service/image"
	"cloth-mini-app/internal/service/item"
	"cloth-mini-app/internal/service/lock"
	"cloth-mini-app/internal/storage/minio"
	"cloth-mini-app/internal/storage/postgresql"
	"fmt"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

// Running application
func Run(config *congig.Config, logger *slog.Logger) {
	logger.Info("starting app...")

	// time.Sleep(time.Second * 15)
	storage, err := postgresql.NewPostgreSQL(config.DB)
	if err != nil {
		logger.Error("failed to init postgresql storage", sl.Err(err))
		fmt.Println(config)
		os.Exit(0)
	}
	_ = storage

	minioClient, err := minio.NewMinioClient(config.Minio)
	if err != nil {
		logger.Error("failed to init minio storage", sl.Err(err))
		os.Exit(1)
	}
	_ = minioClient

	kafkaProducer := kafka.NewProducer(config.Kafka)

	// prepare repositories
	itemRepo := itemRepo.NewItemRepository(logger, storage)
	categoryRepo := categoryRepo.NewCategoryRepository(logger, storage)
	brandRepo := brandRepo.NewBrandRepository(logger, storage)
	imageRepo := imageRepo.NewImageRepository(logger, storage)
	itemImageRepo := itemImageRepo.NewItemImageRepository(logger, storage)
	lockRepo := lockRepo.NewLockRepository(storage)
	outboxRepo := outboxRepo.NewOutboxRepository(logger, storage)

	// facade
	outboxFacade := facade.NewOutboxFacade(storage, logger, outboxRepo, itemImageRepo, brandRepo)

	// prepare services
	lockService := lock.NewLockService(lockRepo)
	itemService := item.NewItemService(logger, itemRepo, imageRepo, itemImageRepo, outboxFacade)
	categoryService := category.NewCategoryService(logger, categoryRepo)
	brandService := brand.NewBrandService(logger, brandRepo)
	imageService := image.NewImageService(logger, minioClient, imageRepo)

	// backgrounds tasks
	backgroundTask := background.NewBackgroundTask(
		logger, minioClient, imageRepo, lockService, outboxRepo, kafkaProducer,
	)
	_ = backgroundTask
	backgroundTask.TempImage.StartDeleteTempImage()
	backgroundTask.Event.StartSendEvent()

	e := echo.New()
	e.Static("/admin/static", "public")

	e.GET("/ping", func(c echo.Context) error {
		c.JSON(200, "pong!")
		return nil
	})

	// prepare handlers
	rest.NewItemHandler(e, itemService)
	rest.NewAdminHandler(e)
	rest.NewCategoryHandler(e, categoryService)
	rest.NewBrandHandler(e, brandService)
	rest.NewImageHandler(e, imageService)

	logger.Info("echo", sl.Err(e.Start(config.Host+":"+config.Port)))
}
