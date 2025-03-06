package repository

import (
	domain "cloth-mini-app/internal/domain/item"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/postgresql"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/Masterminds/squirrel"
)

const (
	limitMax = 100 // max records per query if limit doesn't specified
)

type ItemRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// Get item repository object that represent the item.ItemRepository interface
func NewItemRepository(logger *slog.Logger, db *postgresql.Storage) *ItemRepository {
	return &ItemRepository{
		db:     db.DB,
		logger: logger,
	}
}

// Get items records
func (i *ItemRepository) GetItems(ctx context.Context, params domain.ItemInputData) ([]domain.ItemAPI, error) {
	const op = "repository.item.Items"

	var limit, offset uint64
	if params.Offset != nil {
		offset = uint64(*params.Offset)
	}
	if params.Limit != nil {
		limit = uint64(*params.Limit)
	}

	filter := i.filterItems(params)

	if limit == 0 {
		limit = limitMax
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	q := psql.Select("i.id", "i.name", "i.description", "i.sex", "i.price", "i.discount", "i.outer_link", "i.created_at", "i.updated_at", "c.type", "c.name AS category_name", "b.name").
		From("items i").
		LeftJoin("brand b on i.brand_id = b.id").
		LeftJoin("category c on i.category_id = c.id").
		Limit(limit).
		Offset(offset)

	minPrice, minPriceOk := filter["min_price"]
	maxPrice, maxPriceOk := filter["max_price"]
	if minPriceOk || maxPriceOk {
		if minPrice.(uint) == 0 {
			minPrice = 0
		}
		if maxPrice.(uint) == 0 {
			maxPrice = uint(math.MaxUint)
		}

		q = q.Where(squirrel.Expr("i.price BETWEEN ? AND ?", minPrice, maxPrice))

		delete(filter, "min_price")
		delete(filter, "max_price")
	}

	name, nameOk := filter["i.name"]
	if nameOk {
		q = q.Where(squirrel.Expr("i.name LIKE ?", fmt.Sprintf("%%%v%%", name)))
		delete(filter, "i.name")
	}

	q = q.Where(filter)

	sql, args, _ := q.ToSql()
	// fmt.Println(sql, args, filter)

	rows, err := i.db.Query(sql, args...)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return nil, err
	}
	defer rows.Close()

	var items []domain.ItemAPI
	for rows.Next() {
		var item domain.ItemAPI
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Sex,
			&item.Price,
			&item.Discount,
			&item.OuterLink,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.CategoryType,
			&item.CategoryName,
			&item.BrandName,
		); err != nil {
			i.logger.Error(op, sl.Err(err))

			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Prepare data for where statement
func (i *ItemRepository) filterItems(params domain.ItemInputData) map[string]any {
	filter := make(map[string]any)

	if params.ID != nil {
		filter["i.id"] = *params.ID
	}
	if params.BrandId != nil {
		filter["i.brand_id"] = *params.BrandId
	}
	if params.Name != nil {
		filter["i.name"] = *params.Name
	}
	if params.Sex != nil {
		filter["i.sex"] = *params.Sex
	}
	if params.CategoryId != nil {
		filter["c.id"] = *params.CategoryId
	}

	if params.MinPrice != nil {
		filter["min_price"] = *params.MinPrice
	}
	if params.MaxPrice != nil {
		filter["max_price"] = *params.MaxPrice
	}

	if params.Discount != nil {
		filter["i.discount"] = params.Discount
	}

	return filter
}

// Update item record by ID
func (i *ItemRepository) Update(ctx context.Context, data domain.ItemUpdate) error {
	const op = "repository.item.Update"

	setState := i.updateSetStatements(data)
	if len(setState) == 0 {
		i.logger.Error(fmt.Sprintf("%s: no field to update", op))
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update("items")

	for col, value := range setState {
		psql = psql.Set(col, value)
	}

	sql, args, err := psql.Set("updated_at", time.Now()).Where("id = ?", data.ID).ToSql()
	if err != nil {
		i.logger.Error(op, sl.Err(err))
		return err
	}

	_, err = i.db.Exec(sql, args...)
	if err != nil {
		i.logger.Error(op, sl.Err(err))
		return err
	}

	return nil
}

// Prepare update set statements
// field => value
func (i *ItemRepository) updateSetStatements(item domain.ItemUpdate) map[string]any {
	data := make(map[string]any)

	if item.BrandId != nil {
		data["brand_id"] = *item.BrandId
	}
	if item.Name != nil {
		data["name"] = *item.Name
	}
	if item.Description != nil {
		data["description"] = *item.Description
	}
	if item.CategoryId != nil {
		data["category_id"] = *item.CategoryId
	}
	if item.Sex != nil {
		data["sex"] = *item.Sex
	}
	if item.Discount != nil {
		data["discount"] = *item.Discount
	}
	if item.Price != nil {
		data["price"] = *item.Price
	}
	if item.OuterLink != nil {
		data["outer_link"] = *item.OuterLink
	}

	return data
}

func (i *ItemRepository) GetItemById(ctx context.Context, id int) (domain.ItemAPI, error) {
	const op = "repository.item.ItemById"

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sql, args, err := psql.Select("i.id", "i.name", "i.description", "i.sex", "i.price", "i.discount", "i.outer_link", "i.created_at", "i.updated_at", "c.id as category_id", "c.type", "c.name AS category_name", "b.id as brand_id", "b.name").
		From("items i").
		LeftJoin("brand b on i.brand_id = b.id").
		LeftJoin("category c on i.category_id = c.id").
		Where(squirrel.Expr("i.id = ?", id)).
		ToSql()
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return domain.ItemAPI{}, err
	}

	var item domain.ItemAPI
	err = i.db.QueryRow(sql, args...).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Sex,
		&item.Price,
		&item.Discount,
		&item.OuterLink,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.CategoryId,
		&item.CategoryType,
		&item.CategoryName,
		&item.BrandId,
		&item.BrandName,
	)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return domain.ItemAPI{}, err
	}

	return item, nil
}

func (i *ItemRepository) Create(ctx context.Context, item domain.ItemCreate) error {
	const op = "repository.item.Create"

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("items").
		Columns("brand_id", "name", "description", "sex", "category_id", "price", "discount", "outer_link", "created_at").
		Values(item.BrandId, item.Name, item.Description, item.Sex, item.CategoryId, item.Price, item.Discount, item.OuterLink, time.Now())

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

func (i *ItemRepository) Delete(ctx context.Context, id int) error {
	const op = "repository.item.Delete"

	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Delete("").
		From("items").
		Where("id = ?", id).
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
