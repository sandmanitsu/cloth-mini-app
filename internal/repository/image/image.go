package image

import (
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/postgresql"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
)

type ImageRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewImageRepository(logger *slog.Logger, db *postgresql.Storage) *ImageRepository {
	return &ImageRepository{
		db:     db.DB,
		logger: logger,
	}
}

func (i *ImageRepository) Insert(itemId int, objectId string) error {
	const op = "repository.image.Insert"

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("images").
		Columns("item_id", "object_id").
		Values(itemId, objectId)

	sql, args, err := psql.ToSql()
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return err
	}

	fmt.Println(sql, args)

	_, err = i.db.Exec(sql, args...)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return err
	}

	return nil
}
