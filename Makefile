ifneq (, $(wildcard ./base.env))
	include .env
endif


build:
	@go build -o bin/api -ldflags "-X env.EXT_ENVIRONMENT=dev" cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/api

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@, $(MAKECMDGOALS))

