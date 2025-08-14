# Go API Project

A scalable, idiomatic Go web API project structure with clear separation of concerns, testability, and extensibility.

---

## Project Structure

```
go-api/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/        # HTTP handlers (controllers)
│   │   │   ├── card_handler.go
│   │   │   ├── category_handler.go   # Categories endpoint handler
│   │   │   ├── goodbye_handler.go
│   │   │   ├── health_handler.go
│   │   │   └── user_handler.go
│   │   ├── middleware/      # HTTP middleware
│   │   │   ├── cors.go      # CORS middleware
│   │   │   └── logging.go
│   │   └── routes/          # Route definitions
│   │       └── routes.go
│   ├── domain/
│   │   ├── models/          # Domain models/entities
│   │   │   ├── card.go
│   │   │   ├── response.go  # Standardized API response model
│   │   │   ├── reward_category.go # Reward category model
│   │   │   ├── reward_category_test.go
│   │   │   ├── user.go
│   │   └── interfaces/      # Domain interfaces
│   ├── services/            # Business logic layer
│   │   ├── card_service.go
│   │   ├── card_utils.go
│   │   ├── card_utils_test.go
│   │   ├── category_service.go # Category service
│   │   ├── database_service.go
│   │   ├── goodbye_service.go
│   │   ├── health_service.go
│   │   ├── health_service_test.go
│   │   ├── s3_service.go    # S3 integration service
│   │   └── user_service.go
│   ├── repository/          # Data access layer
│   │   ├── dynamodb/        # DynamoDB specific implementations
│   │   │   └── client.go
│   │   ├── file/            # File-based repository
│   │   │   └── card_repository.go
│   │   └── interfaces/      # Repository interfaces
│   │       ├── card_repository.go
│   │       └── database.go
│   └── config/              # Configuration management
│       └── config.go
├── data/                    # Static data files
│   ├── mock-categories.json # Mock categories data
│   └── mock-recommended-cards.json
├── docs/                    # API documentation (Swagger)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── scripts/                 # Build/deployment scripts
├── test-image.html          # S3 image URL test utility
└── tests/                   # Comprehensive test suite
    ├── unit/                # Unit tests
    ├── integration/         # Integration tests
    ├── e2e/                 # End-to-end tests
    └── helpers/             # Test utilities
```

```
internal/
├── api/          # HTTP handlers and routing (presentation layer)
├── services/     # Business logic (business layer)
├── repository/   # Data access (data layer)
├── domain/       # Business entities/models (domain layer)
└── config/       # Configuration (infrastructure layer)
```

---

## Directory Overview

- **cmd/**: Application entry points (e.g., `main.go` for the server)
- **internal/**: Core application code, protected from external imports
  - **api/handlers/**: HTTP handlers (controllers)
  - **api/middleware/**: HTTP middleware (CORS, logging, etc.)
  - **api/routes/**: Route definitions
  - **domain/models/**: Domain models/entities
  - **domain/interfaces/**: Domain interfaces for abstraction and mocking
  - **services/**: Business logic layer
  - **repository/**: Data access layer (with interfaces and implementations)
  - **config/**: Configuration management
- **data/**: Static/mock data for development and testing
- **docs/**: API documentation (Swagger/OpenAPI)
- **scripts/**: Build and deployment scripts
- **tests/**: Comprehensive test suite (unit, integration, e2e)

---

## Getting Started

### Prerequisites

- Go 1.18+
- (Optional) Docker for running dependencies (e.g., databases)

### Setup

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/go-api.git
   cd go-api
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Run the server:
   ```sh
   go run ./cmd/server/main.go
   ```
4. The API will be available at `http://localhost:10471`.

### Development with Live Reload

For development with automatic reloading when files change:

1. Install Air (live reload tool):

   ```sh
   go install github.com/air-verse/air@latest
   ```

2. Run with live reload:

   ```sh
   air
   ```

3. Make changes to your code and save - the server will automatically rebuild and restart!

---

## Testing

- Run all tests:
  ```sh
  go test ./...
  ```
- Run specific test suites (unit, integration, e2e) from the `tests/` directory as needed.

---

## API Documentation

- Swagger/OpenAPI docs are available in the `docs/` directory.
- To view the docs locally, use a Swagger UI tool or import `swagger.yaml` into [Swagger Editor](https://editor.swagger.io/).

---

## Contribution Guidelines

1. Fork the repository and create your branch from `main`.
2. Follow the existing project structure and naming conventions.
3. Write tests for new features and bug fixes.
4. Run `go fmt` and ensure your code passes linting and tests before submitting a PR.
5. Document any new endpoints or features in the `docs/` directory.

---

## License

[MIT](LICENSE)
