# Makefile for Bankple project

APP_NAME = bankple
CMD_DIR = ./cmd/$(APP_NAME)
BIN_DIR = ./bin
BIN_FILE = $(BIN_DIR)/$(APP_NAME)

.PHONY: all build run test fmt tidy clean

# Default target
all: build

# Build the Go binary
build:
	@echo "🔨 Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_FILE) $(CMD_DIR)/main.go
	@echo "✅ Build complete: $(BIN_FILE)"

# Run the binary
run: build
	@echo "🚀 Running $(APP_NAME)..."
	@$(BIN_FILE)

# Run all tests
test:
	@echo "🧪 Running tests..."
	go test ./...

# Format the code using gofmt
fmt:
	@echo "🧹 Formatting code..."
	gofmt -s -w .

# Tidy up go.mod
tidy:
	go mod tidy

# Clean up the binary
clean:
	@echo "🧼 Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)









































# build:
# 	@go build -o bin/bankple

# run: build
# 	@./bin/bankple

# test:
# 	@go test -v ./...