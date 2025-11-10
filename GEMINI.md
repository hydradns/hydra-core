# Project Overview

This project is a DNS-layer security and privacy gateway called PhantomDNS. It's written in Go and uses a microservices architecture with two main components:

*   **Data Plane**: The core DNS server that handles DNS queries and applies security filtering. It's built using the `miekg/dns` library.
*   **Control Plane**: An administrative API for configuration and monitoring. It's a web service built with the Gin framework.

The two services are designed to be run in Docker containers and communicate with each other via gRPC. The project uses a SQLite database for storage, and GORM as the ORM.

# Architecture

PhantomDNS is composed of two main services: the `dataplane` and the `controlplane`.

## Dataplane

The `dataplane` is the core of PhantomDNS. It is responsible for:

*   **DNS Query Processing**: The `dnsengine` package contains the main logic for handling DNS queries. It listens for incoming queries, processes them, and sends responses.
*   **Blocklist Management**: The `blocklist` package fetches, parses, and stores blocklists from various sources. The `dnsengine` uses these blocklists to filter out malicious or unwanted domains.
*   **Policy Enforcement**: The `policy` package provides a policy engine that allows for more granular control over DNS queries. Policies can be used to block, allow, or redirect queries based on domain names.
*   **Data Persistence**: The `storage` package uses a SQLite database to store blocklists, policies, and query logs.
*   **gRPC Server**: The `grpc` package exposes a gRPC server that the `controlplane` can use to check the health of the `dataplane`.

## Controlplane

The `controlplane` provides a web-based API for managing and monitoring the `dataplane`. It is responsible for:

*   **Admin API**: It uses the Gin framework to expose a RESTful API for managing blocklists, policies, and other settings.
*   **gRPC Client**: It communicates with the `dataplane` via a gRPC client to perform health checks and other administrative tasks.

## Communication

The `controlplane` and `dataplane` communicate with each other using gRPC. The protocol buffer definitions for the gRPC services are located in the `proto` directory. Currently, there is a simple `Health` service defined.

# Building and Running

The project can be built and run using Docker Compose.

## Prerequisites

*   Docker & Docker Compose
*   Go 1.20 or higher (for development)
*   Git

## Running the application

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/lopster568/PhantomDNS.git
    cd PhantomDNS
    ```

2.  **Build and run using Docker Compose:**

    ```sh
    docker-compose up --build
    ```

This will start both the `dataplane` and `controlplane` services.

*   The **Data Plane** (DNS Server) will be listening on port `1053` (UDP and TCP).
*   The **Control Plane** (Admin API) will be listening on port `8086`.

## Building from source

You can also build the services from source:

```sh
# Build the data plane
go build -o bin/dataplane cmd/dataplane/main.go

# Build the control plane
go build -o bin/controlplane cmd/controlplane/main.go
```

# Testing

To run the tests for the project, use the following command:

```sh
go test ./...
```

# Development Conventions

*   The project follows the standard Go project layout.
*   The `internal/` directory contains the private application code, organized by feature (e.g., `blocklist`, `dnsengine`, `policy`).
*   The `cmd/` directory contains the main applications for the `controlplane` and `dataplane` services.
*   Configuration files are located in the `configs/` directory.
*   Dockerfiles are in the `docker/` directory.
*   Protocol Buffer definitions are in the `proto/` directory.