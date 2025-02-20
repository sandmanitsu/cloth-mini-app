package postgresql

import (
	"cloth-mini-app/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

// Create postgresql db instanse
func NewPostgreSQL(cfg config.DB) (*Storage, error) {
	const op = "storage.postgresql.New"

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{DB: db}, nil
}
