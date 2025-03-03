package image

import (
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/postgresql"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
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

func (i *ImageRepository) Insert(itemId int, objectId string) error {
	const op = "repository.image.Insert"

	images, err := i.Images(itemId)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, "failet count existing images"), sl.Err(err))
		return err
	}

	if len(images) >= maxImagesPerItem {
		i.logger.Debug("the number of images per item has reached the maximum", slog.Attr{Key: "itemId", Value: slog.IntValue(itemId)})

		return fmt.Errorf("reached max images per item") // todo. вынести в переменные, что проверить на уровне delivery. В какой слой поместить ???
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("images").
		Columns("item_id", "object_id", "uploaded_at").
		Values(itemId, objectId, time.Now())

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

func (i *ImageRepository) Images(itemId int) ([]string, error) {
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

func (i *ImageRepository) InsertTempImage(objectId string) error {
	const op = "repository.image.InsertTempImage"

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("temp_images").
		Columns("object_id", "uploaded_at").
		Values(objectId, time.Now())

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
