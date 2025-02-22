package brand

import (
	"cloth-mini-app/internal/domain"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/postgresql"
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

func (b *BrandRepository) Brands() ([]domain.Brand, error) {
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
