# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
make dev               # Start server with hot-reload (Air)
make server            # Run server directly
make build             # Build production binary
make test              # Run all tests with coverage
make test-verbose      # Run tests with -failfast
make test-coverage     # Generate HTML coverage report
make repository-mocks  # Regenerate gomock mocks from repository interfaces
make clean             # Clean build artifacts and test cache
```

Run a single test: `go test -run TestFunctionName ./internal/service/token/...`

## Architecture

Clean Architecture with strict inward dependency flow:

**HTTP (api/) → Service (service/) → Repository (repository/) → Model (model/)**

- **api/handler/**: HTTP handlers that parse requests and call services
- **api/middleware/**: JWT auth middleware with CSRF validation
- **api/router.go**: Route definitions (public vs protected)
- **service/**: Business logic, each domain has its own package with interface + implementation
- **repository/interfaces/**: Repository contracts (interfaces)
- **repository/implementations/**: GORM implementations of repository interfaces
- **repository/mock/**: Auto-generated gomock mocks (regenerate with `make repository-mocks`)
- **model/entity/**: GORM database models (UUID PKs, soft deletes)
- **model/request/**: API input DTOs
- **model/response/**: API output DTOs + response wrapper
- **di/container.go**: Wires repositories → services → handlers → router
- **platform/**: Config (env vars) and database initialization (SQLite or PostgreSQL)

Entry point: `cmd/server/main.go` → loads config → builds DI container → starts HTTP server with graceful shutdown.

## Key Conventions

- **TDD workflow**: Write tests first, then implementation (see `internal/service/README.md`)
- **File naming**: `<action>.<layer>.go` (e.g., `get_counter.service.go`, `increment_counter.gorm.go`)
- **Test naming**: `<action>.<layer>_test.go` alongside implementation files
- **Table-driven tests** with `gomock` for repository mocking
- **Context propagation**: All service and repository methods accept `context.Context`
- **Interface-based design**: Services depend on repository interfaces, never concrete types
- **DTOs**: Request models go in, Response models come out — entities stay in the repository layer

## Auth System

- JWT access tokens (15 min) + refresh tokens (7 days) stored in HttpOnly cookies
- CSRF protection via `X-CSRF-Token` header required for state-changing requests
- Bcrypt password hashing
- See `AUTH.md` for full details

## Environment

Copy `env.example` to `.env`. Key vars: `DATABASE_TYPE` (sqlite/postgres), `DATABASE_DSN`, `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`, `CSRF_SECRET`, `DEV_MODE` (enables dev defaults for secrets).
