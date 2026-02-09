.PHONY: all build fmt test rebuild clean

GO_CORE_FILES = cmd/*.go internal/*.go
GO_TEST_FILES = test/*.go

BINARY_NAME = sort

all: build test

build: fmt
	go build -o $(BINARY_NAME) $(GO_CORE_FILES)

fmt:
	go fmt $(GO_CORE_FILES)

test: vet lint
	go test ./...

vet:
	go vet $(GO_CORE_FILES) $(GO_TEST_FILES)

lint:
	-golint $(GO_CORE_FILES) $(GO_TEST_FILES)


rebuild: clean build

clean:
	rm -f $(BINARY_NAME)