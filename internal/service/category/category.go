package category

import (
	domain "cloth-mini-app/internal/domain/category"
	"context"
	"log/slog"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]domain.Category, error)
}

type CategoryService struct {
	logger       *slog.Logger
	categoryRepo CategoryRepository
}

func NewCategoryService(logger *slog.Logger, categoryRepo CategoryRepository) *CategoryService {
	return &CategoryService{
		logger:       logger,
		categoryRepo: categoryRepo,
	}
}

func (c *CategoryService) GetCategories(ctx context.Context) ([]domain.Category, error) {
	return c.categoryRepo.GetCategories(ctx)
}
