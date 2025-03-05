package image

import (
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/postgresql"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
)

var (
	errMaxImages = fmt.Errorf("reached max images per item")
)

const (
	maxImagesPerItem = 4
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

// insert image data to db with SELECT FOR UPDATE
func (i *ImageRepository) Insert(ctx context.Context, itemId int, objectId string) error {
	const op = "repository.image.Insert"

	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, "failet start transaction"), sl.Err(err))
	}
	defer tx.Rollback()

	imagePerItem, err := i.getImagesForUpdate(itemId)
	if err != nil {
		return err
	}

	if imagePerItem >= maxImagesPerItem {
		i.logger.Debug("the number of images per item has reached the maximum", slog.Attr{Key: "itemId", Value: slog.IntValue(itemId)})

		return errMaxImages
	}

	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("images").
		Columns("item_id", "object_id", "uploaded_at").
		Values(itemId, objectId, time.Now()).
		ToSql()
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return err
	}

	_, err = i.db.Exec(sql, args...)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return err
	}

	if err = tx.Commit(); err != nil {
		i.logger.Error(fmt.Sprintf("%s : failed commit transaction", op), sl.Err(err))
	}

	return nil
}

// lock needed rows and return number of images related to provided itemId
func (i *ImageRepository) getImagesForUpdate(itemId int) (int, error) {
	const op = "repository.image.getItemsForUpdate"

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sql, args, err := psql.Select("*").From("images").Where("item_id = ?", itemId).Suffix("for update").ToSql()
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return 0, err
	}

	rows, err := i.db.Query(sql, args...)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return 0, err
	}
	defer rows.Close()

	var imageCnt int
	for rows.Next() {
		if err := rows.Scan(); err != nil {
			i.logger.Error(op, sl.Err(err))

			return 0, err
		}
		imageCnt++
	}

	return imageCnt, nil
}

func (i *ImageRepository) GetImages(itemId int) ([]string, error) {
	const op = "repository.image.Images"

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sql, args, err := psql.Select("object_id").From("images").Where("item_id = ?", itemId).ToSql()
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return nil, err
	}

	rows, err := i.db.Query(sql, args...)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return nil, err
	}
	defer rows.Close()

	var imageIds []string
	for rows.Next() {
		var imageId string
		if err := rows.Scan(&imageId); err != nil {
			i.logger.Error(op, sl.Err(err))

			return nil, err
		}
		imageIds = append(imageIds, imageId)
	}

	return imageIds, nil
}

func (i *ImageRepository) Delete(imageId string) error {
	const op = "repository.image.Delete"

	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Delete("").
		From("images").
		Where("object_id = ?", imageId).
		ToSql()
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return err
	}

	_, err = i.db.Exec(sql, args...)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return err
	}

	return nil
}
