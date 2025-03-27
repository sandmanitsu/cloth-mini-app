run:
	clear
	go run cmd/app/main.go

test:
	go test -v -tags=integrations ./...