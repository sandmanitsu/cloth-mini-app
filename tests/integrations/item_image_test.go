//go:build integrations

package integrations

import (
	"cloth-mini-app/internal/config"
	domain "cloth-mini-app/internal/domain/item"
	itemImageRepo "cloth-mini-app/internal/repository/item_image"
	"cloth-mini-app/internal/storage/postgresql"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ItemImageRepository interface {
	Create(ctx context.Context, item domain.ItemCreate) error
}

type ItemImageIntegrationSuite struct {
	suite.Suite
	db     *sql.DB
	logger *slog.Logger
	repo   ItemImageRepository
}

func NewItemImageSuite() *ItemImageIntegrationSuite {
	return &ItemImageIntegrationSuite{}
}

func (i *ItemImageIntegrationSuite) SetupSuite() {
	config := config.MustLoad("../../.env")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.DBname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	i.db = db

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	i.logger = logger
}

func (i *ItemImageIntegrationSuite) SetupTest() {
	storage := &postgresql.Storage{DB: i.db}

	i.repo = itemImageRepo.NewItemImageRepository(i.logger, storage)
}

func (i *ItemImageIntegrationSuite) TestCreate() {
	ctx := context.Background()

	var images []string
	images = append(images, uuid.NewString())
	images = append(images, uuid.NewString())

	item := domain.ItemCreate{
		BrandId:     1,
		Name:        "create test item",
		Description: "some description...",
		Sex:         1,
		CategoryId:  1,
		Price:       10000,
		Discount:    10,
		OuterLink:   "http:/localhost:8080/",
		Images:      images,
	}

	err := i.repo.Create(ctx, item)
	i.Require().NoError(err)

	dbItem := i.getCreatedItem(item.Name)

	i.Require().Equal(item.Name, dbItem.Name)
	i.Require().Equal(item.Description, dbItem.Description)
	i.Require().Equal(item.Price, dbItem.Price)
	i.Require().Equal(item.BrandId, dbItem.BrandId)
	i.Require().Equal(item.CategoryId, dbItem.CategoryId)
}

func (i *ItemImageIntegrationSuite) getCreatedItem(name string) domain.ItemCreate {
	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(
			"i.id", "i.brand_id", "i.name", "i.description", "i.sex", "i.category_id",
			"i.price", "i.discount", "i.outer_link",
		).
		From("items i").
		ToSql()
	if err != nil {
		log.Fatal(err)
	}

	var itemId int
	var item domain.ItemCreate
	err = i.db.QueryRow(sql, args...).Scan(
		&itemId,
		&item.BrandId,
		&item.Name,
		&item.Description,
		&item.Sex,
		&item.CategoryId,
		&item.Price,
		&item.Discount,
		&item.OuterLink,
	)
	if err != nil {
		log.Fatal(err)
	}

	item.Images = i.getImages(itemId)

	return item
}

func (i *ItemImageIntegrationSuite) getImages(itemId int) []string {
	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("object_id").
		From("images").
		Where("item_id = ?", itemId).
		ToSql()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := i.db.Query(sql, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var image string
		if err := rows.Scan(&image); err != nil {
			log.Fatal(err)
		}

		images = append(images, image)
	}

	return images
}

func (i *ItemImageIntegrationSuite) TearDownTest() {
	_, _ = i.db.Exec("TRUNCATE TABLE temp_images, images, items")
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, NewItemImageSuite())
}
