run:
	clear
	go run cmd/app/main.go

test-integrations:
	docker-compose -f docker-compose.test.yaml -p "integration_tests" up --build --abort-on-container-exit --exit-code-from test

local:
	docker-compose -f docker-compose.yaml up -d