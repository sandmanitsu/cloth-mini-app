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

var (
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

// create item and return itemID and error
func (i *ItemImageRepository) Create(ctx context.Context, item domain.ItemCreate) (uint, error) {
	const op = "repository.item_image.Create"

	var itemID uint
	err := postgresql.WrapTx(ctx, i.db, func(ctx context.Context) error {
		id, err := i.createItem(ctx, item)
		if err != nil {
			return err
		}
		itemID = id

		err = i.createImage(ctx, id, item.Images)
		if err != nil {
			return err
		}

		err = i.deleteFromTempImageTable(ctx, item.Images)
		if err != nil {
			return err
		}

		// err = i.createNotification(ctx, int(itemId), item, eventDomain.EventCreateItem)

		return err
	})

	if err != nil {
		i.logger.Error(op, sl.Err(err))

		return 0, err
	}

	return itemID, nil
}

func (i *ItemImageRepository) createItem(ctx context.Context, item domain.ItemCreate) (uint, error) {
	const op = "repository.item_image.createItem"

	tx, ok := postgresql.TxFromCtx(ctx)
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

	tx, ok := postgresql.TxFromCtx(ctx)
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

	tx, ok := postgresql.TxFromCtx(ctx)
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

func (i *ItemImageRepository) createNotification(ctx context.Context, itemId int, item domain.ItemCreate, event string) error {
	const op = "repository.item_image.createNotification"

	tx, ok := postgresql.TxFromCtx(ctx)
	if !ok {
		i.logger.Error(fmt.Sprintf("%s : failed get transaction from context", op))

		return errGetTransaction
	}

	payload := fmt.Sprintf(`{"item_id":"%d", "brand_id":"%d"}`, itemId, item.BrandId)

	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("outbox").
		Columns("event_type", "payload").
		Values(event, []byte(payload)).
		ToSql()
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return err
	}

	_, err = tx.Exec(sql, args...)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return err
	}

	return fmt.Errorf("some err...")
}
