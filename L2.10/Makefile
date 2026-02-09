.PHONY: all build fmt test rebuild clean

MAIN_PKG = ./cmd/sort
# GO_CORE_FILES = cmd/sort/ internal/
GO_TEST_FILES = test/*.go

BINARY_NAME = bin/sort

all: build

build: fmt
	go build -o $(BINARY_NAME) $(MAIN_PKG)
	./$(BINARY_NAME)

fmt:
	go fmt ./...

test: vet
	go test ./...

vet:
	go vet ./...

lint:
	golint ./...


rebuild: clean build

clean:
	rm -f $(BINARY_NAME)