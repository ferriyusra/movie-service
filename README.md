# Movie Service API

Backend API for a **Movie Reservation System** built with Go, Gin, GORM, and Clean Architecture.

Users can browse movies, view showtimes, and reserve seats. Admins manage movies, theaters, showtimes, and view reports.

## Quick Start

```bash
# Install dependencies
make install-deps

# Copy and configure environment
cp env.example .env

# Run with hot-reload
make dev

# Or run directly
make server
```

The API will be available at `http://localhost:8080`.

## Commands

```bash
make dev              # Start server with hot-reload (air)
make server           # Run server directly
make build            # Build production binary
make test             # Run all tests
make test-coverage    # Run tests with coverage report (generates coverage.html)
make repository-mocks # Regenerate repository mocks after changing interfaces
make clean            # Clean build artifacts
```

## Project Structure

```
cmd/server/main.go          -- Entry point, loads .env, builds DI container, starts Gin
internal/di/container.go     -- Dependency injection: wires repos -> services -> handlers -> routes
```

### Clean Architecture layers (dependencies flow inward only):

```
internal/
├── api/                    # HTTP layer
│   ├── handler/           # Request handlers (bind request, call service, return JSON)
│   ├── middleware/        # JWT auth middleware (Bearer token)
│   └── router.go          # Route registration via SetupRoutes()
│
├── service/               # Business logic layer
│   ├── user/             # User services (register, login, refresh, logout)
│   ├── token/            # JWT token services (access + refresh)
│   └── health/           # Health check services
│
├── repository/           # Data access layer
│   ├── interfaces/       # Repository contracts (*.repository_interface.go)
│   ├── implementations/  # GORM implementations
│   └── mock/            # Auto-generated gomock mocks
│
├── model/                # Data models
│   ├── entity/          # GORM database models (UUID PKs, soft deletes)
│   ├── request/         # API request DTOs
│   └── response/        # API response DTOs (standardized envelope)
│
├── di/                   # Dependency injection container
└── platform/            # Config loading + database initialization (SQLite / PostgreSQL)
```

## API Endpoints

### Public

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login (returns accessToken + refreshToken) |
| POST | `/api/auth/refresh` | Refresh access token |

### Protected (requires `Authorization: Bearer <token>`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/me` | Get current user |
| POST | `/api/auth/logout` | Logout (revokes refresh tokens) |

## Auth System

- **Bearer token** authentication via `Authorization` header
- Access token: **15 min** expiry
- Refresh token: **7 days** expiry, single-use with rotation
- Passwords hashed with **bcrypt**
- Token generation and storage handled in the **service layer** (not handler)

### Auth Flow

```
1. POST /api/auth/register          -- Create account (returns user info only)
2. POST /api/auth/login             -- Returns { accessToken, refreshToken }
3. GET  /api/auth/me                -- Header: Authorization: Bearer <accessToken>
4. POST /api/auth/refresh           -- Body: { "refreshToken": "..." }
5. POST /api/auth/logout            -- Revokes all refresh tokens
```

### Response Format

All responses use a standardized JSON envelope with camelCase keys:

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john@example.com",
      "name": "John Doe"
    },
    "accessToken": "eyJhbG...",
    "refreshToken": "eyJhbG..."
  }
}
```

Error responses:

```json
{
  "success": false,
  "message": "invalid email or password"
}
```

## Environment

Copy `env.example` to `.env`. Key variables:

- `DEV_MODE=true` -- enables dev defaults for JWT secrets
- `DATABASE_TYPE` -- `sqlite` (default) or `postgres`
- `DATABASE_DSN` -- connection string (default `dev.db` for SQLite)
- `JWT_ACCESS_SECRET` / `JWT_REFRESH_SECRET` -- **must** be set in production

## Test-Driven Development (TDD)

This project follows a TDD workflow. All service-layer business logic is covered by unit tests using table-driven patterns and mocked repositories.

### Running Tests

```bash
# Run all tests
make test

# Run specific service tests
go test -run TestLogin ./internal/service/user/...

# Coverage report
make test-coverage
```

### Key Testing Patterns

- **Table-driven tests** -- each test case is a struct in a slice, run via `t.Run()`
- **gomock** -- repository interfaces are mocked, so tests run without a database
- **Context propagation** -- all methods accept `context.Context`, tested with cancellation
- **File convention** -- tests live next to implementation: `login.service.go` + `login.service_test.go`

For the full TDD guide, see [`internal/service/README.md`](internal/service/README.md).
