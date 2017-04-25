TEST = $(shell go list ./... | grep -v '/vendor/')
GOFMT_FILES = $(shell find . -type f -name '*.go' | grep -v vendor)
TREASURY_S3 ?= st-treasury-st-staging
GO_VERSION ?= 1.8.1-alpine
DOCKER_WORKING_DIR := /go/src/github.com/AirHelp/treasury
DOCKER_CMD = docker run --rm -i \
	-e GOOS \
	-v "$(shell pwd)":${DOCKER_WORKING_DIR} \
	-w ${DOCKER_WORKING_DIR} golang:${GO_VERSION}

BUILD_TIME = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_TREE_STATE = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
GIT_IMPORT = github.com/AirHelp/treasury/version
GO_LDFLAGS = -X $(GIT_IMPORT).gitCommit=$(GIT_COMMIT) \
	-X $(GIT_IMPORT).gitTreeState=$(GIT_TREE_STATE) \
	-X $(GIT_IMPORT).buildDate=$(BUILD_TIME)

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
	@GOOS=darwin ${DOCKER_CMD} go build -ldflags "${GO_LDFLAGS}"

dev:
	go build
