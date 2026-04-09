# Clean Architecture Go API (Gin)

A production-ready Go backend built with **Clean Architecture** principles using the **Gin** framework and **GORM** ORM.

This is the backend-only version extracted from [clean-arch-go-vite-react](https://github.com/ferriyusra/clean-arch-go-vite-react), with no frontend embedding or dependencies.

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
cmd/server/main.go          → Entry point, loads config, creates DI container, starts Gin
internal/di/container.go     → Dependency injection: wires repos → services → handlers → routes
```

### Clean Architecture layers (dependencies flow inward only):

```
internal/
├── api/                    # HTTP layer
│   ├── handler/           # Request handlers (bind request, call service, return JSON)
│   ├── middleware/        # JWT auth + CSRF middleware
│   └── router.go          # Route registration via SetupRoutes()
│
├── service/               # Business logic layer
│   ├── user/             # User domain services (register, login, refresh, etc.)
│   ├── counter/          # Counter services
│   ├── csrf/             # CSRF token services
│   ├── token/            # JWT token services
│   ├── health/           # Health check services
│   └── message/          # Message services
│
├── repository/           # Data access layer
│   ├── interfaces/       # Repository contracts (*.repository_interface.go)
│   ├── implementations/  # GORM implementations
│   └── mock/            # Auto-generated gomock mocks
│
├── model/                # Data models
│   ├── entity/          # GORM database models
│   ├── request/         # API request DTOs
│   └── response/        # API response DTOs (standardized envelope)
│
├── di/                   # Dependency injection container
└── platform/            # Config loading + database initialization
```

## API Endpoints

### Public

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/message` | Get message |
| GET | `/api/csrf` | Get CSRF token |
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login |
| POST | `/api/auth/refresh` | Refresh access token (requires CSRF) |

### Protected (requires JWT)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/me` | Get current user |
| POST | `/api/auth/logout` | Logout (requires CSRF) |
| GET | `/api/counter` | Get counter value |
| POST | `/api/counter` | Increment counter (requires CSRF) |

## Auth System

JWT tokens stored in HTTP-only cookies. Access token: 15 min, refresh token: 7 days. CSRF protection via `X-CSRF-Token` header for state-changing requests. See [AUTH.md](AUTH.md) for full details.

## Environment

Copy `env.example` to `.env`. Key variables:

- `DEV_MODE=true` — enables dev defaults for JWT secrets
- `DATABASE_TYPE=sqlite` (default) or `postgres`
- `JWT_ACCESS_SECRET` / `JWT_REFRESH_SECRET` — **must** be set in production
- `CSRF_SECRET` — **must** be set in production

## Test-Driven Development (TDD)

This project follows a TDD workflow — tests are written **before** the implementation. All service-layer business logic is covered by unit tests using table-driven patterns and mocked repositories.

### TDD Cycle

```
1. Write a failing test     → defines the expected behavior
2. Write minimal code       → make the test pass
3. Refactor                 → clean up while keeping tests green
```

### Example: Counter Service

**Step 1 — Write the test first** (`internal/service/counter/get_counter.service_test.go`):

```go
package counter

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ferriyusra/movie-service/internal/repository/mock"
)

func TestGetCounter(t *testing.T) {
	tests := []struct {
		name          string
		mockReturn    int
		mockError     error
		expectedValue int
		expectedError bool
	}{
		{
			name:          "should return counter value successfully",
			mockReturn:    42,
			mockError:     nil,
			expectedValue: 42,
			expectedError: false,
		},
		{
			name:          "should return error when repository fails",
			mockReturn:    0,
			mockError:     errors.New("database connection failed"),
			expectedValue: 0,
			expectedError: true,
		},
		{
			name:          "should return error on context canceled",
			mockReturn:    0,
			mockError:     context.Canceled,
			expectedValue: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockCounterRepository(ctrl)
			mockRepo.EXPECT().
				GetCounter(gomock.Any()).
				Return(tt.mockReturn, tt.mockError).
				Times(1)

			svc := NewCounterService(mockRepo)
			result, err := svc.GetCounter(context.Background())

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected result, got nil")
				} else if result.Value != tt.expectedValue {
					t.Errorf("expected value %d, got %d", tt.expectedValue, result.Value)
				}
			}
		})
	}
}
```

**Step 2 — Implement to make tests pass** (`internal/service/counter/get_counter.service.go`):

```go
package counter

import (
	"context"

	"github.com/ferriyusra/movie-service/internal/model/response"
)

func (s *counterService) GetCounter(ctx context.Context) (*response.GetCounter, error) {
	value, err := s.repo.GetCounter(ctx)
	if err != nil {
		return nil, err
	}
	return &response.GetCounter{
		Value: value,
	}, nil
}
```

**Step 3 — Run and verify**:

```bash
# Run specific service tests
go test ./internal/service/counter -v

# Run all tests
make test
```

### Key Testing Patterns

- **Table-driven tests** — each test case is a struct in a slice, run via `t.Run()`
- **gomock** — repository interfaces are mocked, so tests run without a database
- **Context propagation** — all methods accept `context.Context`, tested with cancellation and deadlines
- **File convention** — tests live next to implementation: `get_counter.service.go` + `get_counter.service_test.go`

For the full TDD guide with step-by-step instructions, see [`internal/service/README.md`](internal/service/README.md).

## Credits

- Originally inspired by [clean-go-vite-react](https://github.com/kamil5b/clean-go-vite-react) by [@kamil5b](https://github.com/kamil5b)
- Extracted from [clean-arch-go-vite-react](https://github.com/ferriyusra/clean-arch-go-vite-react)
