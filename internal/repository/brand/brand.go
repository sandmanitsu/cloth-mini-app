package brand

import (
	domain "cloth-mini-app/internal/domain/brand"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/postgresql"
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
)

type BrandRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewBrandRepository(logger *slog.Logger, db *postgresql.Storage) *BrandRepository {
	return &BrandRepository{
		db:     db.DB,
		logger: logger,
	}
}

func (b *BrandRepository) GetBrands(ctx context.Context) ([]domain.Brand, error) {
	const op = "repository.Brand.Brands"

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, _, err := psql.Select("id", "name").From("Brand").ToSql()
	if err != nil {
		b.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))
	}

	rows, err := b.db.Query(sql)
	if err != nil {
		b.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return nil, err
	}
	defer rows.Close()

	var brands []domain.Brand
	for rows.Next() {
		var brand domain.Brand
		if err := rows.Scan(&brand.ID, &brand.Name); err != nil {
			b.logger.Error(op, sl.Err(err))

			return nil, err
		}
		brands = append(brands, brand)
	}

	return brands, nil
}

func (b *BrandRepository) GetBrand(ctx context.Context, brandId int) (domain.Brand, error) {
	const op = "repository.Brand.GetBrand"

	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "name").
		From("Brand").
		Where("id = ?", brandId).
		ToSql()
	if err != nil {
		b.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))
	}

	var brand domain.Brand
	err = b.db.QueryRow(sql, args...).Scan(&brand.ID, &brand.Name)
	if err != nil {
		b.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return brand, err
	}

	return brand, nil
}
