package repository

import (
	"cloth-mini-app/internal/storage/postgresql"
	"fmt"

	"github.com/Masterminds/squirrel"
)

type ItemRepository struct {
	DB *postgresql.Storage
}

func NewItemRepository(db *postgresql.Storage) *ItemRepository {
	return &ItemRepository{
		DB: db,
	}
}

func (i *ItemRepository) Items() {
	const op = "repository.item.Items"

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, _, _ := psql.Select("*").From("users").ToSql()
	fmt.Println(sql)
}
