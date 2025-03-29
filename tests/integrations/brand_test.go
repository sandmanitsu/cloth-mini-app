package integrations

import (
	"cloth-mini-app/internal/config"
	domain "cloth-mini-app/internal/domain/brand"
	brandRepo "cloth-mini-app/internal/repository/brand"
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

type BrandRepository interface {
	GetBrands(ctx context.Context) ([]domain.Brand, error)
}

type BrandIntegrationSuite struct {
	suite.Suite
	db     *sql.DB
	logger *slog.Logger
	repo   BrandRepository
}

func NewBrandIntegrationSuite() *BrandIntegrationSuite {
	return &BrandIntegrationSuite{}
}

func (b *BrandIntegrationSuite) SetupSuite() {
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

	b.db = db

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	b.logger = logger
}

func (b *BrandIntegrationSuite) SetupTest() {
	storage := &postgresql.Storage{DB: b.db}

	b.repo = brandRepo.NewBrandRepository(b.logger, storage)
}

func (b *BrandIntegrationSuite) TearDownTest() {
	_, _ = b.db.Exec("TRUNCATE TABLE brand")
}

func (b *BrandIntegrationSuite) TestGetBrands() {
	ctx := context.Background()

	brands, err := b.repo.GetBrands(ctx)

	b.Require().NoError(err)
	b.Require().Greater(len(brands), 0)
}

func TestBrandRepoSuite(t *testing.T) {
	suite.Run(t, NewBrandIntegrationSuite())
}
