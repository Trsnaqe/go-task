build:
	@go build -o bin/gotask cmd/main.go

run: build
	@./bin/gotask

test:
	@go test -v ./...

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

swagger:
	@swag init --parseDependency --parseDepth 2 -g cmd/api/api.go 
