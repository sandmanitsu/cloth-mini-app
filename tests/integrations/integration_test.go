//go:build integration

package integrations

import (
	"bytes"
	"cloth-mini-app/internal/config"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
)

const (
	mockItemID = 1 //

	// host = "http://app:8080"
)

var (
	host string
)

type IntegrationSuite struct {
	suite.Suite
	db     *sql.DB
	minio  *minio.Client
	config *config.Config
}

func NewIntegrationSuite() *IntegrationSuite {
	return &IntegrationSuite{}
}

func (i *IntegrationSuite) SetupSuite() {
	i.config = config.MustLoad()

	host = "http://" + i.config.Host + ":" + i.config.Port

	i.getDB()
	i.getMinioClient()
}

func (i *IntegrationSuite) SetupTest() {
	log.Print("migration up")

	err := goose.Up(i.db, "./migrations")
	if err != nil {
		log.Fatal(err)
	}
}

func (i *IntegrationSuite) TearDownTest() {
	log.Print("migration down")

	err := goose.Down(i.db, "./migrations")
	if err != nil {
		log.Fatal(err)
	}
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, NewIntegrationSuite())
}

func (i *IntegrationSuite) getDB() {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		i.config.DB.Host,
		i.config.DB.Port,
		i.config.DB.User,
		i.config.DB.Password,
		i.config.DB.DBname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	i.db = db
}

func (i *IntegrationSuite) getMinioClient() {
	ctx := context.Background()

	client, err := minio.New(i.config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(i.config.Minio.User, i.config.Minio.Password, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	exist, err := client.BucketExists(ctx, i.config.Minio.BucketName)
	if err != nil {
		log.Fatal(err)
	}
	if !exist {
		err := client.MakeBucket(ctx, i.config.Minio.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}

	i.minio = client
}

func (i *IntegrationSuite) getMinioFileId(objectId string) string {
	obj, err := i.minio.GetObject(context.Background(), i.config.Minio.BucketName, objectId, minio.GetObjectOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer obj.Close()

	objInfo, err := obj.Stat()
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, objInfo.Size)
	_, err = obj.Read(buffer)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	return objInfo.Key
}

func (i *IntegrationSuite) putImageToMinio(fileId string, file []byte) {
	reader := bytes.NewReader(file)
	_, err := i.minio.PutObject(context.Background(), i.config.Minio.BucketName, fileId, reader, int64(len(file)), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		log.Fatal(err)
	}
}
