BINARY_NAME=favicreep
BUILD_DIR=bin

.PHONY: all clean build install

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/favicreep

install:
	@echo "Installing $(BINARY_NAME)..."
	@go install ./cmd/favicreep

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
