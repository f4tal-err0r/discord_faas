# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd/
BINARY_NAME := dfaas
EXCLUDE_TEST := "TestContext"

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n s/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /

.PHONY: protobuf
protobuf:
	protoc --go_out=./ --go_opt=paths=source_relative ./proto/*.proto
	protoc --go_out=./pkgs/platform  --go_out=./ --go_opt=Mpkgs/platform/content.proto=github.com/f4tal-err0r/discord_faas/pkgs/platform --go_opt=paths=source_relative ./pkgs/platform/content.proto

.PHONY: protobuf-template
protobuf-template:
	protoc --go_out=./pkgs/platform/templates/golang ./pkgs/platform/content.proto
	protoc --ruby_out=./pkgs/platform/templates/ruby ./pkgs/platform/content.proto

## test: run all tests
.PHONY: test
test: protobuf
	go test -v -race -buildvcs $(shell go list ./... | grep -v 'runtime/') 

## local/test: Exclude tests which can break local environment (like context)
.PHONY: local/test
local/test:
	go test -v -race -buildvcs -tags exclude $(shell go list ./... | grep -v 'runtime/') 

fmt:
	@gofmt -l -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

lint:
	golangci-lint run
	

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover: protobuf
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out $(shell go list ./... | grep -v 'runtime/')
	go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build
build: protobuf
	go build -o=./bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

.PHONY: build/windows
build/windows: protobuf
	env GOOS="windows" GOARCH="amd64" CGO_ENABLED="1" CC="x86_64-w64-mingw32-gcc" go build -o=./bin/${BINARY_NAME}.exe ${MAIN_PACKAGE_PATH}

## run: run the  application
.PHONY: run
run: build
	./bin/${BINARY_NAME}

.PHONY: run/live
run/live: protobuf
	go run github.com/cosmtrek/air@v1.52.0 \
		--build.cmd "make build" --build.bin "./bin/${BINARY_NAME}" --build.delay "100" \
		--build.args_bin "server,start" \
		--build.exclude_dir "runtimes,bin" \
		--misc.clean_on_exit "true"

## production/deploy: deploy the application to production
.PHONY: production/deploy
production/deploy:
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}