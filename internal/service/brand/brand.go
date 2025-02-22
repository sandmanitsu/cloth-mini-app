package brand

import (
	"cloth-mini-app/internal/domain"
	"log/slog"
)

type BrandRepository interface {
	Brands() ([]domain.Brand, error)
}

type BrandService struct {
	logger    *slog.Logger
	BrandRepo BrandRepository
}

func NewBrandService(logger *slog.Logger, BrandRepo BrandRepository) *BrandService {
	return &BrandService{
		logger:    logger,
		BrandRepo: BrandRepo,
	}
}

func (b *BrandService) Brands() ([]domain.Brand, error) {
	return b.BrandRepo.Brands()
}
