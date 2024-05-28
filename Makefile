# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd/
BINARY_NAME := dfaas


# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n s/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build
build:
	go build -o=./bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

.PHONY: build/windows
build/windows:
	env GOOS="windows" GOARCH="amd64" CGO_ENABLED="1" CC="x86_64-w64-mingw32-gcc" go build -o=./bin/${BINARY_NAME}.exe ${MAIN_PACKAGE_PATH}

## run: run the  application
.PHONY: run
run: build
	./bin/${BINARY_NAME}

## run/live: run the application with reloading on file changes
.PHONY: dev/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "./bin/${BINARY_NAME}" --build.delay "100" \
		--misc.clean_on_exit "true"

## production/deploy: deploy the application to production
.PHONY: production/deploy
production/deploy:
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=bin/linux_amd64/${BINARY_NAME} ${MAIN_PACKAGE_PATH}