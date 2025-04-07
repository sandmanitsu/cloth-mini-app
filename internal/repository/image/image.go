package image

import (
	domain "cloth-mini-app/internal/domain/image"
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

		return err
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

		return err
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
		imageCnt++
	}

	return imageCnt, nil
}

func (i *ImageRepository) GetImages(ctx context.Context, itemId int) ([]string, error) {
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

func (i *ImageRepository) Delete(ctx context.Context, imageId string) error {
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

func (i *ImageRepository) InsertTempImage(ctx context.Context, objectId string) error {
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
		if postgresql.IsDuplicateKeyError(err) {
			return nil
		}
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return err
	}

	return nil
}

func (i *ImageRepository) DeleteTempImage(ctx context.Context, deleteFn func([]domain.TempImage) ([]domain.TempImage, error)) error {
	const op = "repository.image.DeleteTempImage"

	images, err := i.getTempImages(ctx)
	if err != nil {
		return err
	}

	images, err = deleteFn(images)
	if err != nil {
		return err
	}

	ids := make([]uint, 0, len(images))
	for _, image := range images {
		ids = append(ids, image.ID)
	}
	if err = i.deleteTempImage(ctx, ids); err != nil {
		return err
	}

	return nil
}

func (i *ImageRepository) getTempImages(ctx context.Context) ([]domain.TempImage, error) {
	const op = "repository.image.getTempImages"

	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "object_id", "uploaded_at").
		From("temp_images").
		Suffix("for update").
		ToSql()
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

	var images []domain.TempImage
	for rows.Next() {
		var image domain.TempImage
		if err := rows.Scan(&image.ID, &image.ObjectId, &image.UploadedAt); err != nil {
			i.logger.Error(op, sl.Err(err))

			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

func (i *ImageRepository) deleteTempImage(ctx context.Context, imageIds []uint) error {
	const op = "repository.image.deleteTempImage"

	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Delete("").
		From("temp_images").
		Where(squirrel.Eq{"id": imageIds}).
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

	i.logger.Info(fmt.Sprintf("%s: delete temp image %d", op, len(imageIds)))

	return nil
}
