BUILD_TIME = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_DISTROS = darwin linux
BUILD_ARCH = amd64 arm64
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_TREE_STATE = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
GIT_IMPORT = github.com/AirHelp/treasury/version
GO_LDFLAGS = -X $(GIT_IMPORT).gitCommit=$(GIT_COMMIT) \
	-X $(GIT_IMPORT).gitTreeState=$(GIT_TREE_STATE) \
	-X $(GIT_IMPORT).buildDate=$(BUILD_TIME)

TREASURY_VERSION?=$(shell awk -F\" '/^const version/ { print $$2; exit }' version/version.go)

default: test

docker-test-build:
	docker-compose -f docker-compose.test.yml build --pull

fmt:
	docker-compose -f docker-compose.test.yml run --rm tests gofmt -s -w .

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	docker-compose -f docker-compose.test.yml run --rm tests go vet -v ./...

test: docker-test-build fmt vet
	docker-compose -f docker-compose.test.yml run --rm tests

testall: test dev
	bats test/bats/tests.bats

build: test
	@rm -fr pkg
	@mkdir pkg
	@for distro in ${BUILD_DISTROS}; do \
		for arch in ${BUILD_ARCH}; do \
			GOOS=$${distro} GOARCH=$${arch} go build -ldflags "${GO_LDFLAGS}" -o pkg/$${distro}/$${arch}/treasury; \
			cd pkg/$${distro}/$${arch}; \
			tar -cjf ../../treasury-$${distro}-$${arch}.tar.bz2 treasury; \
			zip ../../treasury-$${distro}-$${arch}.zip treasury; \
			cd ../../..; \
		done \
	done

_check_deps:
	@which hub >/dev/null || { echo 'No hub cli installed. Exiting...'; exit 1; }
	@which aws >/dev/null || { echo 'No awscli installed. Exiting...'; exit 1; }

_github_release:
	hub release create \
		-a pkg/treasury-darwin-amd64.tar.bz2 \
		-a pkg/treasury-linux-amd64.tar.bz2 \
		-a pkg/treasury-darwin-arm64.tar.bz2 \
		-a pkg/treasury-linux-arm64.tar.bz2 \
		-m ${TREASURY_VERSION} \
		${TREASURY_VERSION}

_s3_release:
	@for distro in ${BUILD_DISTROS}; do \
	    for arch in ${BUILD_ARCH}; do \
			AWS_PROFILE=production aws s3 cp --acl public-read \
				pkg/treasury-$${distro}-$${arch}.tar.bz2 s3://airhelp-devops-binaries/treasury/${TREASURY_VERSION}/treasury-$${distro}-$${arch}.tar.bz2; \
			shasum -a 256 pkg/treasury-$${distro}-$${arch}.tar.bz2; \
			AWS_PROFILE=production aws s3 cp --acl public-read \
				pkg/treasury-$${distro}-$${arch}.zip s3://airhelp-devops-binaries/treasury/${TREASURY_VERSION}/treasury-$${distro}-$${arch}.zip; \
			shasum -a 256 pkg/treasury-$${distro}-$${arch}.zip; \
		done \
	done

_lambda_layers_release: _check_deps
	for env in sta staging development; do \
		aws --profile $${env} lambda publish-layer-version --layer-name treasury-client \
			--description "treasury ${TREASURY_VERSION} layer" \
			--content S3Bucket=airhelp-devops-binaries,S3Key=treasury/${TREASURY_VERSION}/treasury-linux-amd64.zip \
			--compatible-runtimes "python3.8" "go1.x" "nodejs10.x" "nodejs12.x" --output text; \
	done

release: build _check_deps _github_release _s3_release _lambda_layers_release

dev:
	go build
