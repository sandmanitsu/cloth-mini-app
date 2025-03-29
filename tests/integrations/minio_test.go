package integrations

import (
	"cloth-mini-app/internal/config"
	"cloth-mini-app/internal/dto"
	"cloth-mini-app/internal/storage/minio"
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MinioSuite struct {
	suite.Suite
	cl *minio.MinioClient
}

func NewMinioSuite() *MinioSuite {
	return &MinioSuite{}
}

func (m *MinioSuite) SetupSuite() {
	config := config.MustLoad("../../.env")

	cl, err := minio.NewMinioClient(config.Minio)
	if err != nil {
		log.Fatal(err)
	}

	m.cl = cl
}

func (m *MinioSuite) getTestFile() dto.FileDTO {
	data, err := os.ReadFile("fixtures/test_pic.jpg")
	if err != nil {
		log.Fatal(err)
	}

	return dto.FileDTO{
		ID:          uuid.NewString(),
		ContentType: minio.ImageContentType,
		Buffer:      data,
	}
}

func (m *MinioSuite) TestPut() {
	ctx := context.Background()
	file := m.getTestFile()

	err := m.cl.Put(ctx, file)

	m.Require().NoError(err)
}

func (m *MinioSuite) TestGet() {
	ctx := context.Background()
	file := m.getTestFile()

	err := m.cl.Put(ctx, file)
	if err != nil {
		log.Fatal(err)
	}

	minioFile, err := m.cl.Get(ctx, file.ID)

	m.Require().NoError(err)
	m.Require().Equal(file.ID, minioFile.ID)
	m.Require().Equal(file.ContentType, minioFile.ContentType)
	m.Require().Equal(file.Buffer, minioFile.Buffer)
}

func (m *MinioSuite) TestDelete() {
	ctx := context.Background()
	file := m.getTestFile()

	err := m.cl.Put(ctx, file)
	if err != nil {
		log.Fatal(err)
	}

	err = m.cl.Delete(ctx, file.ID)
	if err != nil {
		log.Fatal(err)
	}

	_, err = m.cl.Get(ctx, file.ID)

	m.Require().Equal(err.Error(), "storage.minio.GetImage: The specified key does not exist.")
}

func TestMinioSuite(t *testing.T) {
	suite.Run(t, NewMinioSuite())
}
