package integrations

import (
	"cloth-mini-app/internal/config"
	domain "cloth-mini-app/internal/domain/category"
	categoryRepo "cloth-mini-app/internal/repository/category"
	"cloth-mini-app/internal/storage/postgresql"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]domain.Category, error)
}

type CategoryIntegrationSuite struct {
	suite.Suite
	db     *sql.DB
	logger *slog.Logger
	repo   CategoryRepository
}

func NewCategoryIntegrationSuite() *CategoryIntegrationSuite {
	return &CategoryIntegrationSuite{}
}

func (c *CategoryIntegrationSuite) SetupSuite() {
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

	c.db = db

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	c.logger = logger
}

func (c *CategoryIntegrationSuite) SetupTest() {
	storage := &postgresql.Storage{DB: c.db}

	c.repo = categoryRepo.NewCategoryRepository(c.logger, storage)
}

func (c *CategoryIntegrationSuite) TearDownTest() {
	_, _ = c.db.Exec("TRUNCATE TABLE category")
}

func (c *CategoryIntegrationSuite) TestGetCategorys() {
	ctx := context.Background()

	cat, err := c.repo.GetCategories(ctx)

	c.Require().NoError(err)
	c.Require().Greater(len(cat), 0)
}

func TestCategoryRepoSuite(t *testing.T) {
	suite.Run(t, NewCategoryIntegrationSuite())
}
