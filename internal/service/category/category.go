package category

import (
	"cloth-mini-app/internal/domain"
	"log/slog"
)

type CategoryRepository interface {
	Categories() ([]domain.Category, error)
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

func (c *CategoryService) Categories() ([]domain.Category, error) {
	return c.categoryRepo.Categories()
}
