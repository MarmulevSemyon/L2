# sort test/test_non_sort.txt>> test/sortLinux.txt; ./bin/sort test/test_non_sort.txt>>test/mysort.txt; diff --strip-trailing-cr test/sortLinux.txt test/mysort.txt

APP_NAME := sort
BIN_DIR  := bin
CMD_DIR  := ./cmd/sort

GO      := go
GOFMT   := gofmt
GOLINT  := golint

.PHONY: all fmt vet lint test build run clean

all: fmt vet lint test build

fmt:
	$(GOFMT) -w ./cmd ./internal 

vet:
	$(GO) vet ./...

lint:
	$(GOLINT) ./...

test:
	$(GO) test ./internal/

build:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

run: build
	./$(BIN_DIR)/$(APP_NAME)

clean:
	rm -rf $(BIN_DIR)