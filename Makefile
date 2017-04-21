APP_NAME=treasury
TEST?=$$(go list ./... | grep -v '/vendor/')
GOFMT_FILES?=$$(find . -type f -name '*.go' | grep -v vendor)
TREASURY_S3 := st-treasury-st-staging
GO_VERSION := 1.8.1-alpine
DOCKER_WORKING_DIR := /go/src/github.com/AirHelp/treasury
DOCKER_CMD := docker run --rm -i \
	-e GOOS \
	-v "$$(pwd)":${DOCKER_WORKING_DIR} \
	-w ${DOCKER_WORKING_DIR} golang:${GO_VERSION}

default: test

fmt:
	@echo 'run Go autoformat'
	@${DOCKER_CMD} gofmt -w $(GOFMT_FILES)

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@echo 'run the code static analysis tool'
	@${DOCKER_CMD} go tool vet -all $$(ls -d */ | grep -v vendor)

test: fmt vet
	@echo 'run the unit tests'
	@TREASURY_S3=${TREASURY_S3} \
	${DOCKER_CMD} go test -cover -v $(TEST)

testall: build
	bats test/bats/tests.bats

build: test
	@GOOS=darwin ${DOCKER_CMD} go build

dev:
	GOOS=darwin ${DOCKER_CMD} go build
