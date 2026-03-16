
APP_NAME := telnet
BIN_DIR  := bin
CMD_DIR  := ./cmd/telnet

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

test: build
	$(GO) test ./...

build:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

clean:
	rm -rf $(BIN_DIR)