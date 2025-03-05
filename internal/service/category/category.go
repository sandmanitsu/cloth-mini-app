package category

import (
	domain "cloth-mini-app/internal/domain/category"
	"log/slog"
)

type CategoryRepository interface {
	GetCategories() ([]domain.Category, error)
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

func (c *CategoryService) GetCategories() ([]domain.Category, error) {
	return c.categoryRepo.GetCategories()
}
