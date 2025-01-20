# Variables
BINARY_NAME = dsync
OUTPUT_DIR = bin
MAIN_FILE = ./main.go
INSTALL_DIR = $(GOBIN) # Defaults to $HOME/go/bin if GOBIN is unset

# Targets
.PHONY: build test run i/home/ibergx00/gonstall clean

build:
	@echo "Building the binary..."
	@mkdir -p $(OUTPUT_DIR)
	@go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Binary built at $(OUTPUT_DIR)/$(BINARY_NAME)"

test:
	@echo "Running tests..."
	@go test -v ./...

install: build
	@echo "Installing the binary to $(INSTALL_DIR)..."
	@go install $(MAIN_FILE)
	@echo "Binary installed to $(INSTALL_DIR)"

clean:
	@echo "Cleaning up..."
	@rm -rf $(OUTPUT_DIR)
	@echo "Cleanup complete!"
