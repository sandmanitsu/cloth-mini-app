run:
	clear
	go run cmd/app/main.go

test:
	go test -v -tags=integrations ./...

test-integrations:
	docker-compose -f docker-compose.test.yaml up --build --abort-on-container-exit --exit-code-from test