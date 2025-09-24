# Makefile
run:
	go run main.go

dev:
	air

build:
	go build -o elogika_core main.go

test:
	go test ./...

lint:
	golangci-lint run

docs:
	swag init --parseInternal --parseDependency

typescript:
	go run tools/typescriptify/main.go

migrate:
	go run tools/migrate/main.go

seed:
	go run tools/seeder/main.go
