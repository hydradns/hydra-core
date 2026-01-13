.PHONY: build test lint clean build-docker destroy-docker dev setup fmt vet coverage proto-generate proto-clean

# Proto file directories
PROTO_DIR := proto
PROTO_OUT_DIR := internal/pb
PROTO_VERSION := v1

# Development setup
setup:
	@echo "🔧 Setting up development environment..."
	go mod download
	go mod tidy
	@echo "✅ Development environment ready!"

# Build binaries
build:
	@echo "🏗️ Building PhantomCore..."
	buf lint
	buf generate
	@echo "🔨 Compiling binaries..."
	go build -o bin/controlplane ./cmd/controlplane
	go build -o bin/dataplane ./cmd/dataplane
	@echo "✅ Build completed!"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Run tests with coverage
coverage:
	@echo "📊 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "🔍 Vetting code..."
	go vet ./...

# Lint code
lint:
	@echo "🔧 Linting code..."
	golangci-lint run

# Development with Docker
dev: build-docker
	@echo "🚀 Development environment started!"

build-docker: 
	@echo "🐳 Building Docker containers..."
	docker-compose up -d --build

destroy-docker: 
	@echo "🧹 Destroying Docker containers..."
	docker-compose down

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "✅ Cleanup completed!"

# ============================================================================
# Protocol Buffer Generation
# ============================================================================

# Install protoc dependencies
proto-install:
	@echo "📦 Installing protoc and plugins..."
	@which protoc > /dev/null || (echo "Installing protoc..." && \
		go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
		go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest && \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest)
	@echo "✅ Protoc dependencies installed!"

# Generate protobuf code
proto-generate: proto-install
	@echo "🔨 Generating protobuf code..."
	@mkdir -p $(PROTO_OUT_DIR)/$(PROTO_VERSION)
	@protoc \
		--go_out=$(PROTO_OUT_DIR) \
		--go-grpc_out=$(PROTO_OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		-I proto \
		proto/v1/dataplane.proto
	@echo "✅ Protobuf code generated in $(PROTO_OUT_DIR)!"

# Clean generated protobuf code
proto-clean:
	@echo "🧹 Cleaning generated protobuf code..."
	rm -rf $(PROTO_OUT_DIR)
	@echo "✅ Generated protobuf code cleaned!"

# Validate protobuf definitions
proto-lint:
	@echo "🔍 Linting protobuf definitions..."
	@which buf > /dev/null || (echo "Installing buf..." && go install github.com/bufbuild/buf/cmd/buf@latest)
	@buf lint proto
	@echo "✅ Protobuf definitions valid!"

# ============================================================================
# Help target
# ============================================================================

help:
	@echo "PhantomDNS Development Commands:"
	@echo ""
	@echo "Build & Development:"
	@echo "  setup        - Setup development environment"
	@echo "  build        - Build binaries"
	@echo "  test         - Run tests"
	@echo "  coverage     - Run tests with coverage report"
	@echo "  dev          - Start development environment"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Lint code"
	@echo ""
	@echo "Protocol Buffers:"
	@echo "  proto-install   - Install protoc and plugins"
	@echo "  proto-generate  - Generate protobuf code"
	@echo "  proto-lint      - Validate protobuf definitions"
	@echo "  proto-clean     - Clean generated code"
	@echo ""
	@echo "Docker:"
	@echo "  build-docker    - Build and start Docker containers"
	@echo "  destroy-docker  - Stop and remove Docker containers"
	@echo ""
	@echo "Utilities:"
	@echo "  clean        - Clean build artifacts"
	@echo "  help         - Show this help message"