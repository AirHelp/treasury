APP_NAME=treasury
TEST?=$$(go list ./... | grep -v '/vendor/')
GOFMT_FILES?=$$(find . -type f -name '*.go' | grep -v vendor)

default: test

fmt:
	@echo 'run Go autoformat'
	@gofmt -w $(GOFMT_FILES)

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@echo 'run the code static analysis tool'
	@go tool vet -all $$(ls -d */ | grep -v vendor)

test: fmt vet
	@echo 'run the unit tests'
	@TREASURY_S3=st-treasury-st-staging \
	go test -cover -v $(TEST)

testall: build
	bats test/bats/tests.bats

build: test
	go build

dev:
	go build
