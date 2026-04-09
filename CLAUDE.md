# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
make install-deps      # go mod tidy + download
make dev               # Start server with hot-reload (Air), sets DEV_MODE=true
make server            # Run server directly with DEV_MODE=true
make build             # Build production binary to ./bin/
make test              # Run all tests with -v -cover -race
make test-verbose      # Run tests with -failfast
make test-coverage     # Generate HTML coverage report (coverage.html)
make repository-mocks  # Regenerate gomock mocks from repository interfaces
make clean             # Clean build artifacts and test cache
```

Run a single test: `go test -run TestFunctionName ./internal/service/token/...`

Go module: `github.com/ferriyusra/movie-service` (differs from directory name).

## Architecture

Clean Architecture with strict inward dependency flow:

**HTTP (api/) → Service (service/) → Repository (repository/) → Model (model/)**

- **api/handler/**: HTTP handlers that parse requests and call services
- **api/middleware/**: JWT auth middleware with CSRF validation; sets `user_id`, `user_email`, and `claims` context keys on authenticated requests
- **api/router.go**: Route definitions — public routes, auth routes (some with CSRF), and protected routes behind `AuthMiddleware`
- **service/**: Business logic, each domain has its own package with interface + implementation. Token service is standalone (no repository dependency, configured via `TokenConfig` struct)
- **repository/interfaces/**: Repository contracts (`*_interface.go` files)
- **repository/implementations/**: GORM implementations, each in its own subdomain package
- **repository/mock/**: Auto-generated gomock mocks (regenerate with `make repository-mocks`)
- **model/entity/**: GORM database models (UUID PKs, soft deletes)
- **model/request/**: API input DTOs
- **model/response/**: API output DTOs + `response.Err()`/`response.OK()` wrapper
- **di/container.go**: Wires repositories → services → handlers → router
- **platform/**: Config (env vars via `godotenv`) and database initialization (SQLite or PostgreSQL)

Entry point: `cmd/server/main.go` → loads `.env` → builds DI container → starts HTTP server with graceful shutdown.

## Key Conventions

- **TDD workflow**: Write tests first, then implementation (see `internal/service/README.md`)
- **File naming**: `<action>.<layer>.go` (e.g., `get_counter.service.go`, `increment_counter.gorm.go`)
- **Test naming**: `<action>.<layer>_test.go` alongside implementation files
- **Table-driven tests** with `gomock` for repository mocking
- **Context propagation**: All service and repository methods accept `context.Context`
- **Interface-based design**: Services depend on repository interfaces, never concrete types
- **DTOs**: Request models go in, Response models come out — entities stay in the repository layer

## Auth System

- JWT access tokens (15 min) + refresh tokens (7 days) stored in HttpOnly cookies (`access_token`, `refresh_token`)
- CSRF protection via `X-CSRF-Token` header required for POST/PUT/PATCH/DELETE requests
- Bcrypt password hashing
- Middleware context keys: `user_id` (string UUID), `user_email`, `claims` (`*token.TokenClaims`)
- Helper functions in middleware: `GetUserIDFromContext(c)`, `GetEmailFromContext(c)`, `GetClaimsFromContext(c)`

## Environment

Copy `env.example` to `.env`. Key vars:
- `DEV_MODE` — when true, uses fallback dev secrets for JWT/CSRF (never use in production)
- `DATABASE_TYPE` — `sqlite` (default) or `postgres`; note: not in env.example, set via `DATABASE_TYPE` env var
- `DATABASE_DSN` — defaults to `dev.db` for SQLite
- `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`, `CSRF_SECRET` — required in production (`DEV_MODE=false`)
- `ALLOWED_ORIGINS` — comma-separated, defaults to `http://localhost:5173`
