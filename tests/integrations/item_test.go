package integrations

import (
	"cloth-mini-app/internal/config"
	domain "cloth-mini-app/internal/domain/item"
	itemRepo "cloth-mini-app/internal/repository/item"
	"cloth-mini-app/internal/storage/postgresql"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/suite"
)

type ItemRepository interface {
	// Fetch items from db
	GetItems(ctx context.Context, params domain.ItemInputData) ([]domain.ItemAPI, error)
	// Returning item by id
	GetItemById(ctx context.Context, id int) (domain.ItemAPI, error)
	// Update item record
	Update(ctx context.Context, data domain.ItemUpdate) error
	// Delete item
	Delete(ctx context.Context, id int) error
}

type ItemIntegrationSuite struct {
	suite.Suite
	db     *sql.DB
	logger *slog.Logger
	repo   ItemRepository
}

func NewItemSuite() *ItemIntegrationSuite {
	return &ItemIntegrationSuite{}
}

func (i *ItemIntegrationSuite) SetupSuite() {
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

func (i *ItemIntegrationSuite) SetupTest() {
	storage := &postgresql.Storage{DB: i.db}

	i.repo = itemRepo.NewItemRepository(i.logger, storage)
}

func (i *ItemIntegrationSuite) TearDownTest() {
	_, _ = i.db.Exec("TRUNCATE TABLE items")
}

func TestItemRepoSuite(t *testing.T) {
	suite.Run(t, NewItemSuite())
}

func (i *ItemIntegrationSuite) TestDelete() {
	ctx := context.Background()

	testItem := domain.ItemCreate{
		BrandId:     1,
		Name:        "test delete item",
		Description: "some description...",
		Sex:         1,
		CategoryId:  1,
		Price:       10000,
		Discount:    10,
		OuterLink:   "http:/localhost:8080/",
		Images:      nil,
	}

	id := i.createItem(testItem)

	i.repo.Delete(ctx, int(id))

	_, err := i.repo.GetItemById(ctx, int(id))

	i.Require().Error(err)
}

func (i *ItemIntegrationSuite) TestGetItemById() {
	ctx := context.Background()

	testItem := domain.ItemCreate{
		BrandId:     1,
		Name:        "test get item by it",
		Description: "some description...",
		Sex:         1,
		CategoryId:  1,
		Price:       10000,
		Discount:    10,
		OuterLink:   "http:/localhost:8080/",
		Images:      nil,
	}

	id := i.createItem(testItem)

	dbItem, err := i.repo.GetItemById(ctx, int(id))

	i.Require().NoError(err)

	i.Require().Equal(id, dbItem.ID)
}

func (i *ItemIntegrationSuite) TestUpdateItem() {
	ctx := context.Background()

	testItem := domain.ItemCreate{
		BrandId:     1,
		Name:        "test update item",
		Description: "some description...",
		Sex:         1,
		CategoryId:  1,
		Price:       10000,
		Discount:    10,
		OuterLink:   "http:/localhost:8080/",
		Images:      nil,
	}

	i.createItem(testItem)

	newName := "Updated"
	var newPrice uint = 500

	err := i.repo.Update(ctx, domain.ItemUpdate{
		Name:  &newName,
		Price: &newPrice,
	})

	i.Require().NoError(err)

	dbItems, err := i.repo.GetItems(ctx, domain.ItemInputData{
		Name: &newName,
	})

	i.Require().NoError(err)

	i.Require().Equal(newName, dbItems[0].Name)
	i.Require().Equal(newPrice, dbItems[0].Price)
}

func (i *ItemIntegrationSuite) TestGetItem() {
	ctx := context.Background()

	testItems := []domain.ItemCreate{
		{
			BrandId:     1,
			Name:        "test get item 1",
			Description: "some description...",
			Sex:         1,
			CategoryId:  1,
			Price:       10000,
			Discount:    10,
			OuterLink:   "http:/localhost:8080/",
			Images:      nil,
		},
		{
			BrandId:     1,
			Name:        "test get item 2",
			Description: "some description...",
			Sex:         1,
			CategoryId:  1,
			Price:       10000,
			Discount:    10,
			OuterLink:   "http:/localhost:8080/",
			Images:      nil,
		},
		{
			BrandId:     1,
			Name:        "test get item 3",
			Description: "some description...",
			Sex:         1,
			CategoryId:  1,
			Price:       10000,
			Discount:    10,
			OuterLink:   "http:/localhost:8080/",
			Images:      nil,
		},
	}

	for _, testItem := range testItems {
		i.createItem(testItem)
	}

	dbItems, err := i.repo.GetItems(ctx, domain.ItemInputData{})

	i.Require().NoError(err)

	i.Require().Equal(len(testItems), len(dbItems))
}

func (i *ItemIntegrationSuite) TestGetItemByName() {
	ctx := context.Background()

	item := domain.ItemCreate{
		BrandId:     1,
		Name:        "create test item",
		Description: "some description...",
		Sex:         1,
		CategoryId:  1,
		Price:       10000,
		Discount:    10,
		OuterLink:   "http:/localhost:8080/",
		Images:      nil,
	}

	i.createItem(item)

	items, err := i.repo.GetItems(ctx, domain.ItemInputData{
		Name: &item.Name,
	})

	i.Require().NoError(err)

	dbItem := items[0]

	i.Require().Equal(item.Name, dbItem.Name)
	i.Require().Equal(item.Description, dbItem.Description)
}

func (i *ItemIntegrationSuite) createItem(item domain.ItemCreate) uint {
	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("items").
		Columns("brand_id", "name", "description", "sex", "category_id", "price", "discount", "outer_link", "created_at").
		Values(item.BrandId, item.Name, item.Description, item.Sex, item.CategoryId, item.Price, item.Discount, item.OuterLink, time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		log.Fatal(err)
	}

	var itemId uint
	err = i.db.QueryRow(sql, args...).Scan(&itemId)
	if err != nil {
		log.Fatal(err)
	}

	return itemId
}
