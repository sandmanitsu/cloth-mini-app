run:
	clear
	HOST=localhost \
	PORT=8081 \
	ENV=dev \
	DBHOST=localhost \
	USER=admin \
	PASSWORD=admin \
	DBNAME=storage \
	DBPORT=5430 \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING=postgres://$USER:$PASSWORD@$DBHOST:$DBPORT/$DBNAME \
	GOOSE_MIGRATION_DIR=./migrations \
	MINIO_ENDPOINT=localhost:9000 \
	MINIO_BUCKET_NAME=image-bucket \
	MINIO_ROOT_USER=admin \
	MINIO_ROOT_PASSWORD=minio123 \
	KAFKA_BROKER=localhost:9094 \
	KAFKA_TOPIC=notifications \
	go run cmd/app/main.go

migration_up:
	clear
	HOST=localhost \
	PORT=8081 \
	ENV=dev \
	DBHOST=localhost \
	USER=admin \
	PASSWORD=admin \
	DBNAME=storage \
	DBPORT=5430 \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING=postgres://$USER:$PASSWORD@$DBHOST:$DBPORT/$DBNAME \
	GOOSE_MIGRATION_DIR=./migrations \
	MINIO_ENDPOINT=localhost:9000 \
	MINIO_BUCKET_NAME=image-bucket \
	MINIO_ROOT_USER=admin \
	MINIO_ROOT_PASSWORD=minio123 \
	KAFKA_BROKER=localhost:9094 \
	KAFKA_TOPIC=notifications \
	go run cmd/migration/migration.go

test-integrations:
	docker-compose -f docker-compose.test.yaml -p "integration_tests" up --build --abort-on-container-exit --exit-code-from test

local:
	docker-compose -f docker-compose.yaml up -d