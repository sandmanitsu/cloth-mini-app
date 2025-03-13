package repository

import (
	domain "cloth-mini-app/internal/domain/item"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/postgresql"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
)

type txKey string

var (
	// tx key in context
	ctxTxKey = txKey("tx")

	errGetTransaction = fmt.Errorf("error: getting transaction from context")
)

type ItemImageRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewItemImageRepository(logger *slog.Logger, db *postgresql.Storage) *ItemImageRepository {
	return &ItemImageRepository{
		db:     db.DB,
		logger: logger,
	}
}

func (i *ItemImageRepository) Create(ctx context.Context, item domain.ItemCreate) error {
	const op = "repository.item_image.Create"

	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, "failet start transaction"), sl.Err(err))
	}
	defer tx.Rollback()

	ctx = context.WithValue(ctx, ctxTxKey, tx)

	itemId, err := i.createItem(ctx, item)
	if err != nil {
		return err
	}

	err = i.createImage(ctx, itemId, item.Images)
	if err != nil {
		return err
	}

	err = i.deleteFromTempImageTable(ctx, item.Images)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		i.logger.Error(fmt.Sprintf("%s : failed commit transaction", op), sl.Err(err))

		return err
	}

	return nil
}

func (i *ItemImageRepository) txFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(ctxTxKey).(*sql.Tx)
	if !ok {
		return nil, false
	}

	return tx, true
}

func (i *ItemImageRepository) createItem(ctx context.Context, item domain.ItemCreate) (uint, error) {
	const op = "repository.item_image.createItem"

	tx, ok := i.txFromCtx(ctx)
	if !ok {
		i.logger.Error(fmt.Sprintf("%s : failed get transaction from context", op))

		return 0, errGetTransaction
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("items").
		Columns("brand_id", "name", "description", "sex", "category_id", "price", "discount", "outer_link", "created_at").
		Values(item.BrandId, item.Name, item.Description, item.Sex, item.CategoryId, item.Price, item.Discount, item.OuterLink, time.Now()).
		Suffix("RETURNING id")

	sql, args, err := psql.ToSql()
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return 0, err
	}

	var itemId uint
	err = tx.QueryRow(sql, args...).Scan(&itemId)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return 0, err
	}

	return itemId, nil
}

type Image struct {
	ItemId   uint
	FileId   string
	UploadAt time.Time
}

func (i *ItemImageRepository) createImage(ctx context.Context, itemId uint, images []string) error {
	const op = "repository.item_image.createImage"

	tx, ok := i.txFromCtx(ctx)
	if !ok {
		i.logger.Error(fmt.Sprintf("%s : failed get transaction from context", op))

		return errGetTransaction
	}

	sql, _, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("images").
		Columns("item_id", "object_id", "uploaded_at").
		Values("", "", "").
		ToSql()
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return err
	}

	args := make([]Image, 0, len(images))
	for _, image := range images {
		args = append(args, Image{
			ItemId:   itemId,
			FileId:   image,
			UploadAt: time.Now(),
		})
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return err
	}
	defer stmt.Close()

	for _, arg := range args {
		_, err := stmt.Exec(arg.ItemId, arg.FileId, arg.UploadAt)
		if err != nil {
			i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

			return err
		}
	}

	return nil
}

func (i *ItemImageRepository) deleteFromTempImageTable(ctx context.Context, imageIds []string) error {
	const op = "repository.item_image.deleteFromTempImageTable"

	tx, ok := i.txFromCtx(ctx)
	if !ok {
		i.logger.Error(fmt.Sprintf("%s : failed get transaction from context", op))

		return errGetTransaction
	}

	sql, _, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Delete("").
		From("temp_images").
		Where("object_id = ?").
		ToSql()
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return err
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return err
	}

	for _, id := range imageIds {
		_, err := stmt.Exec(id)
		if err != nil {
			i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

			return err
		}
	}

	return nil
}
