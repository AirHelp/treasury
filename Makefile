GOFMT_FILES = $(shell find . -type f -name '*.go' | grep -v vendor)
TREASURY_S3 ?= ah-dev-treasury-development
DOCKER_TEST_IMAGE := airhelp/treasury-test
DOCKER_WORKING_DIR := /go/src/github.com/AirHelp/treasury
DOCKER_CMD = docker run --rm -i \
	-e GOOS \
	-v "$(shell pwd)":${DOCKER_WORKING_DIR} \
	-w ${DOCKER_WORKING_DIR} ${DOCKER_TEST_IMAGE}

BUILD_TIME = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_DISTROS = darwin linux
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_TREE_STATE = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
GIT_IMPORT = github.com/AirHelp/treasury/version
GO_LDFLAGS = -X $(GIT_IMPORT).gitCommit=$(GIT_COMMIT) \
	-X $(GIT_IMPORT).gitTreeState=$(GIT_TREE_STATE) \
	-X $(GIT_IMPORT).buildDate=$(BUILD_TIME)

TREASURY_VERSION?=$(shell awk -F\" '/^const version/ { print $$2; exit }' version/version.go)

default: test

docker-test-build:
	docker build -t $(DOCKER_TEST_IMAGE) -f Dockerfile-test .

fmt:
	@echo 'run Go autoformat'
	@${DOCKER_CMD} gofmt -w $(GOFMT_FILES)

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@echo 'run the code static analysis tool'
	@${DOCKER_CMD} go tool vet -all $$(ls -d */ | grep -v vendor)

test: docker-test-build fmt vet
	@echo 'run the unit tests'
	@TREASURY_S3=${TREASURY_S3} \
	${DOCKER_CMD} go test -cover -v ./...

testall: test dev
	bats test/bats/tests.bats

build: test
	@rm -fr pkg
	@mkdir pkg
	@for distro in ${BUILD_DISTROS}; do \
		GOOS=$${distro} ${DOCKER_CMD} go build -ldflags "${GO_LDFLAGS}" -o pkg/$${distro}/treasury; \
		cd pkg/$${distro}; \
		tar -cjf ../treasury-$${distro}-amd64.tar.bz2 treasury; \
		cd ../..; \
	done

release: build
	@which hub >/dev/null || { echo 'No hub cli installed. Exiting...'; exit 1; }
	hub release create \
		-a pkg/treasury-darwin-amd64.tar.bz2 \
		-a pkg/treasury-linux-amd64.tar.bz2 \
		-m ${TREASURY_VERSION} \
		${TREASURY_VERSION}
	@for distro in ${BUILD_DISTROS}; do \
		AWS_PROFILE=production aws s3 cp --acl public-read \
			pkg/treasury-$${distro}-amd64.tar.bz2 s3://airhelp-devops-binaries/treasury/${TREASURY_VERSION}/treasury-$${distro}-amd64.tar.bz2; \
		shasum -a 256 pkg/treasury-$${distro}-amd64.tar.bz2; \
	done
dev:
	go build
