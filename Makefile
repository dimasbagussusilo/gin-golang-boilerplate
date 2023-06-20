NAME=gin-golang-boilerplate
VERSION=0.0.1

.PHONY: generatedocs1
## Create initiation migration
generatedocs1:
	@rm -r docs/v1
	@mkdir docs/v1 && swag init -g routers/v1.go --output docs/v1/

.PHONY: build
## build: Compile the packages.
build:
	@go build -o $(NAME)

.PHONY: run
## run: Build and Run in development mode.
run: build
	@./$(NAME)

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f $(NAME)

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

.PHONY: dev
## dev: Running development
dev:
	nodemon --exec go run main.go --signal SIGTERM

.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
