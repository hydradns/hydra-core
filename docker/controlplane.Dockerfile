# docker/controlplane.Dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install build dependencies including protobuf compiler
RUN apk add --no-cache git protobuf-dev protoc

COPY go.mod go.sum ./
RUN go mod download

# Install protobuf Go plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .

# Ensure proto files are up to date (regenerate if needed)
RUN if [ -f proto/health.proto ]; then \
    protoc --go_out=. --go-grpc_out=. proto/health.proto; \
    fi

RUN go build -o controlplane ./cmd/controlplane

# Final runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/controlplane .

# Run as non-root (we’ll map port 53 later)
USER 1000:1000

CMD ["./controlplane"]
