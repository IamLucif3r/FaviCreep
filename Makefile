BINARY_NAME=favicreep
BUILD_DIR=bin

.PHONY: all clean build install

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

install:
	@echo "Installing $(BINARY_NAME)..."
	@go install main.go

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
