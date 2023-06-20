NAME=gin-golang-boilerplate
VERSION=0.0.1

## This file is not git ignored, please do not put real credentials here
DB_URL=postgres://postgres:password@localhost:5432/testing?sslmode=disable

all: help
# help: show this help message.
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

# generatedocs1: generate docs v1.
generatedocs1:
	@rm -r docs/v1
	@mkdir docs/v1 && swag init -g routers/v1.go --output docs/v1/

# migrationcreate: create initialize migration.
migrationcreate:
	migrate create -ext sql -dir database/migrations -seq init_schema

# migrateup: start migration up.
migrateup:
	migrate -path database/migrations -database "$(DB_URL)" -verbose up

# migratedown: start migration down.
migratedown:
	migrate -path database/migrations -database "$(DB_URL)" -verbose down

# build: compile the packages to build binary file.
build:
	@go build -o $(NAME)

# run: build and run project.
run: build
	@./$(NAME)

# clean: Clean project and previous builds.
clean:
	@rm -f $(NAME)

# deps: download all dependencies project.
deps:
	@go mod download

# dev: running development with nodemon.
dev:
	nodemon --exec go run main.go --signal SIGTERM

.PHONY: help generatedocs1 migrationcreate migratedown build run clean deps dev
