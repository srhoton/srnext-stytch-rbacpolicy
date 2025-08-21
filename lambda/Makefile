.PHONY: all build test coverage clean fmt lint vet install-tools

# Build variables
BINARY_NAME=bootstrap
BUILD_DIR=bin
LAMBDA_DIR=cmd/lambda
GOOS=linux
GOARCH=arm64
GO_VERSION=1.24

# Test variables
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html
COVERAGE_THRESHOLD=80

all: clean fmt vet lint test build

install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

build:
	@echo "Building Lambda function for ARM64..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -tags lambda.norpc -o $(BUILD_DIR)/$(BINARY_NAME) ./$(LAMBDA_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test:
	@echo "Running tests..."
	go test -v -race ./...

coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@echo "Generating coverage report..."
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@echo "Checking coverage threshold..."
	@coverage=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$coverage%"; \
	if [ $$(echo "$$coverage < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "Coverage $$coverage% is below threshold $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	else \
		echo "Coverage $$coverage% meets threshold $(COVERAGE_THRESHOLD)%"; \
	fi

fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

vet:
	@echo "Running go vet..."
	go vet ./...

lint:
	@echo "Running linter..."
	golangci-lint run --timeout=5m

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

# Development helpers
run-local:
	@echo "Running Lambda locally (requires SAM CLI)..."
	sam local start-api

package:
	@echo "Creating deployment package..."
	mkdir -p $(BUILD_DIR)
	$(MAKE) build
	cd $(BUILD_DIR) && zip -r lambda-deployment.zip $(BINARY_NAME)
	@echo "Deployment package created: $(BUILD_DIR)/lambda-deployment.zip"

# Docker build for consistent environment
docker-build:
	@echo "Building in Docker container..."
	docker run --rm -v "$$(pwd)":/workspace -w /workspace public.ecr.aws/lambda/provided:al2023 \
		sh -c "dnf install -y golang && make build"

# Dependency management
deps:
	@echo "Downloading dependencies..."
	go mod download

tidy:
	@echo "Tidying dependencies..."
	go mod tidy

verify:
	@echo "Verifying dependencies..."
	go mod verify

# Help target
help:
	@echo "Available targets:"
	@echo "  all            - Run fmt, vet, lint, test, and build"
	@echo "  build          - Build the Lambda binary for ARM64"
	@echo "  test           - Run all tests"
	@echo "  coverage       - Run tests with coverage analysis"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run golangci-lint"
	@echo "  clean          - Remove build artifacts"
	@echo "  package        - Create deployment package"
	@echo "  docker-build   - Build in Docker container"
	@echo "  deps           - Download dependencies"
	@echo "  tidy           - Tidy dependencies"
	@echo "  verify         - Verify dependencies"
	@echo "  install-tools  - Install development tools"
	@echo "  help           - Show this help message"