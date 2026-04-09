# 🎬 Product Requirements Document (PRD)
## Movie Reservation System

---

| Field         | Detail                          |
|---------------|---------------------------------|
| **Version**   | 1.0.0                           |
| **Status**    | Draft                           |
| **Author**    | Engineering Team                |
| **Tech Stack**| Go (Golang) · Next.js · PostgreSQL |
| **Last Updated** | 2025                         |

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Goals & Objectives](#2-goals--objectives)
3. [User Roles & Permissions](#3-user-roles--permissions)
4. [Functional Requirements](#4-functional-requirements)
5. [Non-Functional Requirements](#5-non-functional-requirements)
6. [Data Model & Entity Relationships](#6-data-model--entity-relationships)
7. [API Design](#7-api-design)
8. [Business Logic & Edge Cases](#8-business-logic--edge-cases)
9. [Frontend Pages & User Flows](#9-frontend-pages--user-flows)
10. [Tech Stack & Architecture](#10-tech-stack--architecture)
11. [Security Considerations](#11-security-considerations)
12. [Reporting & Analytics](#12-reporting--analytics)
13. [Out of Scope](#13-out-of-scope)
14. [Future Enhancements](#14-future-enhancements)
15. [Glossary](#15-glossary)

---

## 1. Project Overview

The **Movie Reservation System** is a full-stack web application that enables users to browse currently showing movies, view available showtimes, and reserve specific seats. The platform supports two user roles — **regular users** and **admins** — with different capabilities scoped to each role.

The backend is built with **Go (Gin framework)**, the frontend with **Next.js (App Router)**, and data is persisted in **PostgreSQL**. Authentication is token-based using **JWT**.

---

## 2. Goals & Objectives

### Primary Goals
- Provide a seamless movie browsing and seat reservation experience for end users.
- Enable admins to manage the complete movie catalog, showtimes, and theaters.
- Prevent overbooking through atomic seat-locking mechanisms.
- Give admins visibility into reservations, capacity, and revenue.

### Learning Objectives (Engineering)
- Implement complex relational data models with proper foreign key constraints.
- Handle concurrent seat reservation using database-level locking (`SELECT FOR UPDATE`).
- Design role-based access control (RBAC) with JWT middleware.
- Write complex SQL queries for reporting and aggregation.

---

## 3. User Roles & Permissions

### Role Definitions

| Permission                        | Guest | User | Admin |
|-----------------------------------|:-----:|:----:|:-----:|
| Browse movies                     | ✅    | ✅   | ✅    |
| View showtimes                    | ✅    | ✅   | ✅    |
| Sign up / Log in                  | ✅    | —    | —     |
| Reserve seats                     | ❌    | ✅   | ✅    |
| View own reservations             | ❌    | ✅   | ✅    |
| Cancel own upcoming reservations  | ❌    | ✅   | ✅    |
| Add / Edit / Delete movies        | ❌    | ❌   | ✅    |
| Add / Edit / Delete showtimes     | ❌    | ❌   | ✅    |
| Manage theaters & seats           | ❌    | ❌   | ✅    |
| Manage genres                     | ❌    | ❌   | ✅    |
| View all reservations             | ❌    | ❌   | ✅    |
| View revenue & capacity reports   | ❌    | ❌   | ✅    |
| Promote users to admin            | ❌    | ❌   | ✅    |

### Initial Admin
The system ships with one **seed admin account** bootstrapped on first startup:
```
Email:    admin@cinema.com
Password: Admin@123  (must be changed after first login)
Role:     admin
```

---

## 4. Functional Requirements

### 4.1 Authentication & Authorization

#### FR-AUTH-01 — User Registration
- Users can register with: `name`, `email`, `password`.
- Email must be unique across the system.
- Password must be minimum 8 characters, containing at least one uppercase letter and one number.
- Passwords are stored as bcrypt hashes (cost factor ≥ 12).
- On success, return a JWT access token and refresh token.

#### FR-AUTH-02 — User Login
- Users log in with `email` and `password`.
- On success, return a JWT access token (24h expiry) and refresh token (7 days expiry).
- On failure (wrong credentials), return a generic `401 Unauthorized` — do not reveal whether the email exists.

#### FR-AUTH-03 — Token Refresh
- Clients can exchange a valid refresh token for a new access token.
- Refresh tokens are single-use (rotated on every refresh).

#### FR-AUTH-04 — Role-Based Authorization Middleware
- Every protected route validates the JWT signature and expiry.
- Admin-only routes additionally verify `role = "admin"`.
- Attempting to access admin routes as a regular user returns `403 Forbidden`.

#### FR-AUTH-05 — Promote User to Admin
- An admin can promote any regular user to admin by user ID.
- An admin cannot demote themselves.

---

### 4.2 Movie Management (Admin Only)

#### FR-MOVIE-01 — Create Movie
Fields required: `title`, `description`, `poster_url`, `duration_min`, `release_date`, `language`, `genre_ids[]`.

#### FR-MOVIE-02 — Update Movie
Admin can update any movie field including reassigning genres.

#### FR-MOVIE-03 — Delete Movie
- Soft delete: sets `deleted_at` timestamp, preserves historical reservation data.
- Cannot delete a movie that has future showtimes with active reservations.

#### FR-MOVIE-04 — List Movies (Public)
- Supports filtering by: `genre`, `date` (has showtime on that date), `title` (search).
- Supports pagination: `page`, `limit` (default 20).
- Returns movies with their genres (not showtime details).

#### FR-MOVIE-05 — Get Movie Detail (Public)
Returns full movie info including genres. Showtimes for a queried date can be optionally included.

#### FR-GENRE-01 — Genre Management (Admin)
- CRUD operations on genres (name must be unique).
- A genre can only be deleted if no movies are tagged with it.

---

### 4.3 Theater & Seat Management (Admin Only)

#### FR-THEATER-01 — Create Theater
Fields: `name`, `location`, `rows` (A–Z), `seats_per_row` (1–30).  
On creation, seats are auto-generated: e.g., `A1, A2, ... Z30`.

#### FR-THEATER-02 — Update Theater
Name and location only. Seat layout cannot be changed after creation if the theater has past reservations.

#### FR-SEAT-01 — Seat Generation
Each seat has: `id`, `theater_id`, `row` (letter), `number` (integer), `label` (e.g., `B7`), `type` (standard/vip/accessible).

---

### 4.4 Showtime Management (Admin Only)

#### FR-SHOW-01 — Create Showtime
Fields: `movie_id`, `theater_id`, `start_time` (UTC datetime), `price` (decimal).  
`end_time` is auto-calculated: `start_time + movie.duration_min + 15min` (cleaning buffer).

#### FR-SHOW-02 — Conflict Detection
Two showtimes in the same theater cannot overlap. The system validates:
```
new.start_time < existing.end_time AND new.end_time > existing.start_time
```
If conflict detected, return `409 Conflict` with conflicting showtime details.

#### FR-SHOW-03 — Update / Delete Showtime
- Cannot update a showtime that has confirmed reservations.
- Deleting a showtime with reservations requires admin confirmation; all reservations are auto-cancelled and users are notified (in-app notification or email — see Future Enhancements).

#### FR-SHOW-04 — List Showtimes by Date
Returns all showtimes for a given `date` (local date mapped to UTC range), grouped by movie.  
Response includes `available_seats` count.

---

### 4.5 Seat Reservation (Authenticated Users)

#### FR-RES-01 — Browse Available Seats
Given a `showtime_id`, return the theater's full seat map with status per seat:
- `available` — not reserved
- `reserved` — taken by another user
- `selected` — in the user's current cart (session-scoped, TTL 10 min)

#### FR-RES-02 — Create Reservation
- User submits: `showtime_id`, `seat_ids[]` (1–10 seats per reservation).
- System performs within a **single DB transaction**:
  1. `SELECT … FOR UPDATE` on all requested seats to lock rows.
  2. Check none of the seats is already reserved for this showtime.
  3. Create `Reservation` record with status `confirmed`.
  4. Create `ReservationSeat` records linking each seat.
  5. Calculate `total_amount = seat_count × showtime.price`.
- If any seat is taken, the transaction rolls back and returns `409 Conflict` with the conflicting seat labels.
- Returns reservation summary with a unique `booking_reference` (e.g., `CINE-20250601-00042`).

#### FR-RES-03 — View My Reservations
Returns paginated list of user's reservations with: movie title, showtime date/time, theater name, seat labels, total amount, status.  
Supports filter by status: `confirmed`, `cancelled`.

#### FR-RES-04 — Cancel Reservation
- User can cancel only their own reservations.
- Cancellation is only allowed if `showtime.start_time > now + 1 hour` (configurable buffer).
- On cancel: status changes to `cancelled`; seats are released back to available.
- Returns `403` if reservation belongs to another user.
- Returns `422 Unprocessable Entity` if cancellation window has passed.

---

### 4.6 Admin — Reservation Management

#### FR-ADMIN-RES-01 — List All Reservations
Paginated list of all reservations across all users.  
Supports filters: `user_id`, `movie_id`, `showtime_id`, `status`, `date_from`, `date_to`.

#### FR-ADMIN-RES-02 — Cancel Any Reservation
Admin can cancel any reservation regardless of timing window.

---

### 4.7 Reporting (Admin Only)

#### FR-REPORT-01 — Revenue Report
Total revenue grouped by: movie, theater, date range.  
Response: `{ movie_title, showtime_count, tickets_sold, total_revenue }`.

#### FR-REPORT-02 — Capacity Report
Per showtime: `total_seats`, `seats_reserved`, `seats_available`, `occupancy_rate (%)`.

#### FR-REPORT-03 — Dashboard Summary
Single endpoint returning:
- Total movies, showtimes, users, active reservations today.
- Revenue this week vs. last week.
- Top 5 movies by tickets sold.
- Top 5 movies by revenue.

---

## 5. Non-Functional Requirements

| Category         | Requirement |
|------------------|-------------|
| **Performance**  | API response < 200ms for read operations under normal load (< 100 concurrent users) |
| **Concurrency**  | Seat reservation must be safe under concurrent requests (no overbooking) |
| **Security**     | Passwords bcrypt-hashed; JWT signed with HS256; HTTPS only in production |
| **Reliability**  | Database transactions with rollback on any failure |
| **Scalability**  | Stateless backend — horizontally scalable behind a load balancer |
| **Pagination**   | All list endpoints paginated (default `limit=20`, max `limit=100`) |
| **Validation**   | All inputs validated server-side; descriptive error messages with field names |
| **Logging**      | Structured JSON logs for all requests (method, path, status, latency, user_id) |
| **CORS**         | Configurable allowed origins; frontend origin whitelisted |
| **Timezone**     | All datetimes stored in UTC; clients responsible for local conversion |

---

## 6. Data Model & Entity Relationships

### 6.1 Entity Relationship Overview

```
┌──────────┐      ┌──────────────┐      ┌──────────┐
│  users   │ 1──N │ reservations │ N──1 │showtimes │
└──────────┘      └──────┬───────┘      └─────┬────┘
                         │                    │
                    1     │                N  N│  1
                         N                    │
               ┌──────────────────┐      ┌────┴─────┐      ┌────────┐
               │ reservation_seats│      │  movies  │ N──N │ genres │
               └──────────┬───────┘      └──────────┘      └────────┘
                          │
                     N    │    1
                    ┌─────┴──────┐
                    │   seats    │ N──1 ┌──────────┐
                    └────────────┘      │ theaters │
                                        └──────────┘
```

### 6.2 Table Schemas

#### `users`
```sql
id            SERIAL PRIMARY KEY
name          VARCHAR(255) NOT NULL
email         VARCHAR(255) UNIQUE NOT NULL
password      VARCHAR(255) NOT NULL          -- bcrypt hash
role          VARCHAR(20) DEFAULT 'user'     -- 'user' | 'admin'
created_at    TIMESTAMPTZ DEFAULT NOW()
updated_at    TIMESTAMPTZ DEFAULT NOW()
deleted_at    TIMESTAMPTZ                    -- soft delete
```

#### `genres`
```sql
id            SERIAL PRIMARY KEY
name          VARCHAR(100) UNIQUE NOT NULL
created_at    TIMESTAMPTZ DEFAULT NOW()
```

#### `movies`
```sql
id            SERIAL PRIMARY KEY
title         VARCHAR(255) NOT NULL
description   TEXT
poster_url    VARCHAR(500)
duration_min  INTEGER NOT NULL
language      VARCHAR(50) DEFAULT 'English'
release_date  DATE
rating        DECIMAL(3,1) DEFAULT 0.0
created_at    TIMESTAMPTZ DEFAULT NOW()
updated_at    TIMESTAMPTZ DEFAULT NOW()
deleted_at    TIMESTAMPTZ                    -- soft delete
```

#### `movie_genres` (join table)
```sql
movie_id      INTEGER REFERENCES movies(id) ON DELETE CASCADE
genre_id      INTEGER REFERENCES genres(id) ON DELETE RESTRICT
PRIMARY KEY (movie_id, genre_id)
```

#### `theaters`
```sql
id            SERIAL PRIMARY KEY
name          VARCHAR(255) NOT NULL
location      VARCHAR(255)
total_seats   INTEGER NOT NULL
rows          INTEGER NOT NULL               -- number of rows
seats_per_row INTEGER NOT NULL
created_at    TIMESTAMPTZ DEFAULT NOW()
updated_at    TIMESTAMPTZ DEFAULT NOW()
```

#### `seats`
```sql
id            SERIAL PRIMARY KEY
theater_id    INTEGER NOT NULL REFERENCES theaters(id)
row           CHAR(1) NOT NULL               -- 'A', 'B', ...
number        INTEGER NOT NULL               -- 1, 2, 3, ...
label         VARCHAR(10) NOT NULL           -- 'A1', 'B12', ...
type          VARCHAR(20) DEFAULT 'standard' -- 'standard' | 'vip' | 'accessible'
UNIQUE (theater_id, row, number)
```

#### `showtimes`
```sql
id            SERIAL PRIMARY KEY
movie_id      INTEGER NOT NULL REFERENCES movies(id)
theater_id    INTEGER NOT NULL REFERENCES theaters(id)
start_time    TIMESTAMPTZ NOT NULL
end_time      TIMESTAMPTZ NOT NULL           -- auto-calculated
price         DECIMAL(10,2) NOT NULL
created_at    TIMESTAMPTZ DEFAULT NOW()
updated_at    TIMESTAMPTZ DEFAULT NOW()
deleted_at    TIMESTAMPTZ
```

#### `reservations`
```sql
id                SERIAL PRIMARY KEY
booking_reference VARCHAR(30) UNIQUE NOT NULL  -- e.g. CINE-20250601-00042
user_id           INTEGER NOT NULL REFERENCES users(id)
showtime_id       INTEGER NOT NULL REFERENCES showtimes(id)
status            VARCHAR(20) DEFAULT 'confirmed' -- 'confirmed' | 'cancelled'
total_amount      DECIMAL(10,2) NOT NULL
cancelled_at      TIMESTAMPTZ
created_at        TIMESTAMPTZ DEFAULT NOW()
updated_at        TIMESTAMPTZ DEFAULT NOW()
```

#### `reservation_seats`
```sql
id              SERIAL PRIMARY KEY
reservation_id  INTEGER NOT NULL REFERENCES reservations(id)
seat_id         INTEGER NOT NULL REFERENCES seats(id)
UNIQUE (seat_id, showtime_id)  -- enforced via application + DB constraint
```

> **Critical Index:** Add a partial unique index to prevent double-booking:
> ```sql
> CREATE UNIQUE INDEX idx_seat_showtime_active
>   ON reservation_seats (seat_id, showtime_id)
>   WHERE status = 'confirmed';
> -- where showtime_id is denormalized or joined from reservations
> ```

---

## 7. API Design

### Base URL
```
https://api.yourdomain.com/api/v1
```

### Authentication Header
```
Authorization: Bearer <access_token>
```

### Standard Response Envelope
```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 150
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "SEAT_ALREADY_RESERVED",
    "message": "Seats B3, B4 are already reserved.",
    "fields": { "seat_ids": ["B3 is taken", "B4 is taken"] }
  }
}
```

---

### 7.1 Auth Endpoints

| Method | Path                    | Auth   | Description                   |
|--------|-------------------------|--------|-------------------------------|
| POST   | `/auth/register`        | None   | Register new user             |
| POST   | `/auth/login`           | None   | Login, receive tokens         |
| POST   | `/auth/refresh`         | None   | Refresh access token          |
| GET    | `/auth/me`              | User   | Get current user profile      |

### 7.2 Movie Endpoints

| Method | Path                    | Auth   | Description                   |
|--------|-------------------------|--------|-------------------------------|
| GET    | `/movies`               | None   | List movies (with filters)    |
| GET    | `/movies/:id`           | None   | Movie detail                  |
| POST   | `/movies`               | Admin  | Create movie                  |
| PUT    | `/movies/:id`           | Admin  | Update movie                  |
| DELETE | `/movies/:id`           | Admin  | Soft delete movie             |
| GET    | `/genres`               | None   | List all genres               |
| POST   | `/genres`               | Admin  | Create genre                  |
| DELETE | `/genres/:id`           | Admin  | Delete genre                  |

### 7.3 Theater Endpoints

| Method | Path                    | Auth   | Description                   |
|--------|-------------------------|--------|-------------------------------|
| GET    | `/theaters`             | None   | List all theaters             |
| GET    | `/theaters/:id`         | None   | Theater detail with seats     |
| POST   | `/theaters`             | Admin  | Create theater + auto-gen seats |
| PUT    | `/theaters/:id`         | Admin  | Update theater info           |

### 7.4 Showtime Endpoints

| Method | Path                        | Auth   | Description                        |
|--------|-----------------------------|--------|------------------------------------|
| GET    | `/showtimes?date=YYYY-MM-DD`| None   | Showtimes for a date               |
| GET    | `/showtimes/:id`            | None   | Showtime detail                    |
| GET    | `/showtimes/:id/seats`      | None   | Seat map with availability         |
| POST   | `/showtimes`                | Admin  | Create showtime                    |
| PUT    | `/showtimes/:id`            | Admin  | Update showtime                    |
| DELETE | `/showtimes/:id`            | Admin  | Delete showtime                    |

### 7.5 Reservation Endpoints

| Method | Path                        | Auth   | Description                        |
|--------|-----------------------------|--------|------------------------------------|
| GET    | `/reservations`             | User   | My reservations (paginated)        |
| GET    | `/reservations/:id`         | User   | Reservation detail                 |
| POST   | `/reservations`             | User   | Create reservation                 |
| DELETE | `/reservations/:id`         | User   | Cancel reservation                 |

### 7.6 Admin Endpoints

| Method | Path                                | Auth  | Description                   |
|--------|-------------------------------------|-------|-------------------------------|
| GET    | `/admin/reservations`               | Admin | All reservations (filtered)   |
| DELETE | `/admin/reservations/:id`           | Admin | Force cancel any reservation  |
| GET    | `/admin/reports/revenue`            | Admin | Revenue report                |
| GET    | `/admin/reports/capacity`           | Admin | Capacity report               |
| GET    | `/admin/reports/dashboard`          | Admin | Dashboard summary             |
| GET    | `/admin/users`                      | Admin | List all users                |
| PATCH  | `/admin/users/:id/promote`          | Admin | Promote user to admin         |

---

## 8. Business Logic & Edge Cases

### 8.1 Overbooking Prevention

The core mechanism uses **PostgreSQL row-level locking** inside a transaction:

```sql
BEGIN;

-- Step 1: Lock the seats for this showtime
SELECT rs.seat_id
FROM reservation_seats rs
JOIN reservations r ON r.id = rs.reservation_id
WHERE rs.seat_id = ANY($1)
  AND r.showtime_id = $2
  AND r.status = 'confirmed'
FOR UPDATE;

-- Step 2: Check if any seats came back (they're taken)
-- If rows returned → ROLLBACK with 409

-- Step 3: If clear → INSERT reservation + reservation_seats

COMMIT;
```

Additionally, a **database-level partial unique index** ensures no two `confirmed` reservations share the same `(seat_id, showtime_id)`:
```sql
CREATE UNIQUE INDEX idx_no_double_booking
ON reservation_seats(seat_id, (
  SELECT showtime_id FROM reservations WHERE id = reservation_id
))
WHERE (
  SELECT status FROM reservations WHERE id = reservation_id
) = 'confirmed';
```

### 8.2 Showtime Scheduling Rules

- A theater can only host **one movie at a time**.
- Scheduling checks: `new.start_time < existing.end_time AND new.end_time > existing.start_time`.
- Showtimes must start at least **30 minutes from the current time** when created.
- End time = `start_time + duration_min + 15` (buffer for cleaning/ads).

### 8.3 Cancellation Rules

| Condition                               | Outcome        |
|-----------------------------------------|----------------|
| `showtime.start_time > now + 1 hour`    | ✅ Allowed      |
| `showtime.start_time ≤ now + 1 hour`    | ❌ Rejected (422) |
| Reservation status = `cancelled`        | ❌ Already cancelled (400) |
| Reservation belongs to another user     | ❌ Forbidden (403) |

### 8.4 Seat Map Rendering

The seat availability endpoint returns a structured response for UI rendering:

```json
{
  "theater": { "name": "Hall A", "rows": 8, "seats_per_row": 10 },
  "seats": [
    { "id": 1, "label": "A1", "row": "A", "number": 1, "type": "standard", "status": "available" },
    { "id": 2, "label": "A2", "row": "A", "number": 2, "type": "standard", "status": "reserved" },
    ...
  ]
}
```

Clients render this as a grid, coloring by status:
- 🟩 `available` — Green
- 🟥 `reserved` — Red / Gray
- 🟦 `selected` — Blue (client-side state)

### 8.5 Booking Reference Generation

```
Format: CINE-{YYYYMMDD}-{zero-padded-id}
Example: CINE-20250601-00042
```

Generated server-side after insertion using the reservation's auto-increment ID.

---

## 9. Frontend Pages & User Flows

### 9.1 Public Pages

| Route                    | Description                                        |
|--------------------------|----------------------------------------------------|
| `/`                      | Home: featured movies, currently showing          |
| `/movies`                | Movie listing with genre filter + search          |
| `/movies/[id]`           | Movie detail: synopsis, genres, showtimes by date |
| `/showtimes/[id]/seats`  | Seat picker for a specific showtime               |
| `/login`                 | Login form                                        |
| `/register`              | Registration form                                 |

### 9.2 User Pages (Authenticated)

| Route                    | Description                                        |
|--------------------------|----------------------------------------------------|
| `/reservations`          | My reservation history                            |
| `/reservations/[id]`     | Reservation detail + ticket view                  |
| `/profile`               | User profile settings                             |

### 9.3 Admin Pages

| Route                        | Description                                 |
|------------------------------|---------------------------------------------|
| `/admin`                     | Dashboard: stats overview                   |
| `/admin/movies`              | Movie list + CRUD                           |
| `/admin/movies/new`          | Create movie form                           |
| `/admin/movies/[id]/edit`    | Edit movie form                             |
| `/admin/theaters`            | Theater management                          |
| `/admin/showtimes`           | Showtime calendar/list + CRUD               |
| `/admin/reservations`        | All reservations with filters               |
| `/admin/reports`             | Revenue & capacity reports                  |
| `/admin/users`               | User list + promote to admin                |

### 9.4 Key User Flows

#### Flow 1: Book a Seat
```
Home → Browse Movies → Select Movie
  → View Showtimes for Date → Select Showtime
    → Seat Picker → Select Seats → Confirm
      → Reservation Confirmed (booking reference shown)
```

#### Flow 2: Cancel Reservation
```
My Reservations → Select Reservation
  → Click Cancel → Confirm Dialog
    → Success: Seats released
```

#### Flow 3: Admin Add Movie
```
Admin Dashboard → Movies → New Movie
  → Fill form (title, description, poster, genres, duration)
    → Save → Movie appears in listing
```

---

## 10. Tech Stack & Architecture

### Backend

| Layer           | Technology                                |
|-----------------|-------------------------------------------|
| Language        | Go 1.25+                                  |
| Framework       | Gin v1.12                                 |
| ORM             | GORM v1.31                                |
| Database        | PostgreSQL (gorm/driver/postgres) or SQLite (glebarez/sqlite) |
| Auth            | JWT (golang-jwt/jwt/v5)                   |
| Password Hash   | bcrypt (golang.org/x/crypto)              |
| Config          | godotenv                                  |
| Validation      | go-playground/validator/v10 (via Gin)     |
| UUID            | google/uuid                               |
| Mocking         | golang/mock (gomock)                      |
| CORS            | gin-contrib/cors                          |

### Backend Architecture (Clean / Layered)

Strict inward dependency flow: **HTTP (api/) → Service (service/) → Repository (repository/) → Model (model/)**

File naming convention: `<action>.<layer>.go` (e.g., `create_movie.service.go`, `find_by_id.gorm.go`)
Test naming convention: `<action>.<layer>_test.go` alongside implementation files.

```
cmd/
└── server/
    └── main.go                          ← Entrypoint, loads .env, builds DI, graceful shutdown

internal/
├── api/
│   ├── handler/                         ← HTTP request/response layer (1 file per domain)
│   │   ├── user.go                      ✅ Exists — register, login, logout, refresh, getMe
│   │   ├── health.go                    ✅ Exists — health check
│   │   ├── counter.go                   ✅ Exists — demo feature (can remove later)
│   │   ├── message.go                   ✅ Exists — demo feature (can remove later)
│   │   ├── movie.go                     🔲 TODO — CRUD movies (admin) + list/detail (public)
│   │   ├── genre.go                     🔲 TODO — CRUD genres (admin) + list (public)
│   │   ├── theater.go                   🔲 TODO — CRUD theaters (admin) + list/detail (public)
│   │   ├── showtime.go                  🔲 TODO — CRUD showtimes (admin) + list/seats (public)
│   │   ├── reservation.go              🔲 TODO — create/cancel/list reservations (user)
│   │   └── admin.go                     🔲 TODO — admin reservations, reports, user promotion
│   ├── middleware/                       ← Auth, RBAC, CSRF, logger, CORS
│   │   ├── auth.go                      ✅ Exists — JWT auth + optional auth + CSRF validation
│   │   └── rbac.go                      🔲 TODO — AdminOnly middleware (role check)
│   └── router.go                        ✅ Exists — route registration (needs new domain routes)
│
├── di/
│   └── container.go                     ✅ Exists — DI wiring (needs new domain wiring)
│
├── model/
│   ├── entity/                          ← GORM database models (UUID PKs, soft deletes)
│   │   ├── user.go                      ✅ Exists — id, email, password, name, role, timestamps
│   │   ├── refresh_token.go             ✅ Exists — id, user_id, token, expires_at
│   │   ├── counter.go                   ✅ Exists — demo (can remove later)
│   │   ├── message.go                   ✅ Exists — demo (can remove later)
│   │   ├── movie.go                     🔲 TODO — title, description, poster_url, duration_min, etc.
│   │   ├── genre.go                     🔲 TODO — name (unique)
│   │   ├── movie_genre.go               🔲 TODO — join table (movie_id, genre_id)
│   │   ├── theater.go                   🔲 TODO — name, location, rows, seats_per_row, total_seats
│   │   ├── seat.go                      🔲 TODO — theater_id, row, number, label, type
│   │   ├── showtime.go                  🔲 TODO — movie_id, theater_id, start_time, end_time, price
│   │   ├── reservation.go              🔲 TODO — booking_reference, user_id, showtime_id, status, total_amount
│   │   └── reservation_seat.go          🔲 TODO — reservation_id, seat_id
│   ├── request/                         ← API input DTOs
│   │   ├── user.go                      ✅ Exists — RegisterUserRequest, LoginRequest
│   │   ├── movie.go                     🔲 TODO — CreateMovieRequest, UpdateMovieRequest
│   │   ├── genre.go                     🔲 TODO — CreateGenreRequest
│   │   ├── theater.go                   🔲 TODO — CreateTheaterRequest, UpdateTheaterRequest
│   │   ├── showtime.go                  🔲 TODO — CreateShowtimeRequest, UpdateShowtimeRequest
│   │   └── reservation.go              🔲 TODO — CreateReservationRequest
│   └── response/                        ← API output DTOs + response envelope
│       ├── wrapper.go                   ✅ Exists — APIResponse, Meta, OK/OKWithMeta/Err helpers
│       ├── common.go                    ✅ Exists — CommonIDResponse
│       ├── user.go                      ✅ Exists — GetUser, LoginResponse, RegisterResponse
│       ├── health.go                    ✅ Exists
│       ├── counter.go                   ✅ Exists — demo
│       ├── message.go                   ✅ Exists — demo
│       ├── movie.go                     🔲 TODO — MovieResponse, MovieListResponse
│       ├── genre.go                     🔲 TODO — GenreResponse
│       ├── theater.go                   🔲 TODO — TheaterResponse, SeatMapResponse
│       ├── showtime.go                  🔲 TODO — ShowtimeResponse
│       ├── reservation.go              🔲 TODO — ReservationResponse, ReservationListResponse
│       └── report.go                    🔲 TODO — RevenueReport, CapacityReport, DashboardSummary
│
├── platform/                            ← Infrastructure concerns
│   ├── config.go                        ✅ Exists — env var loading (AuthConfig, ServerConfig, DatabaseConfig)
│   ├── config_test.go                   ✅ Exists
│   └── database.go                      ✅ Exists — GORM init (SQLite or PostgreSQL)
│
├── repository/
│   ├── interfaces/                      ← Repository contracts (1 interface file per domain)
│   │   ├── user.repository_interface.go           ✅ Exists
│   │   ├── refresh_token.repository_interface.go  ✅ Exists
│   │   ├── counter.repository_interface.go        ✅ Exists — demo
│   │   ├── message.repository_interface.go        ✅ Exists — demo
│   │   ├── movie.repository_interface.go          🔲 TODO — Create, Update, Delete, FindByID, FindAll
│   │   ├── genre.repository_interface.go          🔲 TODO — Create, Delete, FindAll, FindByID
│   │   ├── theater.repository_interface.go        🔲 TODO — Create, Update, FindByID, FindAll
│   │   ├── seat.repository_interface.go           🔲 TODO — CreateBatch, FindByTheaterID
│   │   ├── showtime.repository_interface.go       🔲 TODO — Create, Update, Delete, FindByID, FindByDate, CheckConflict
│   │   ├── reservation.repository_interface.go    🔲 TODO — Create, Cancel, FindByID, FindByUserID, FindAll
│   │   └── reservation_seat.repository_interface.go 🔲 TODO — CreateBatch, FindByShowtimeID, CheckAvailability
│   ├── implementations/                 ← GORM implementations (1 package per domain, 1 file per method)
│   │   ├── user/
│   │   │   ├── user.gorm.go             ✅ — struct + NewGORMUserRepository
│   │   │   ├── create.gorm.go           ✅
│   │   │   ├── find_by_id.gorm.go       ✅
│   │   │   ├── find_by_email.gorm.go    ✅
│   │   │   ├── update.gorm.go           ✅
│   │   │   └── delete.gorm.go           ✅
│   │   ├── refresh_token/
│   │   │   ├── refresh_token.gorm.go    ✅
│   │   │   ├── create.gorm.go           ✅
│   │   │   ├── find_by_token.gorm.go    ✅
│   │   │   ├── delete_by_token.gorm.go  ✅
│   │   │   └── delete_by_user_id.gorm.go ✅
│   │   ├── counter/                     ✅ — demo
│   │   ├── message/                     ✅ — demo
│   │   ├── movie/                       🔲 TODO
│   │   │   ├── movie.gorm.go            🔲 — struct + constructor
│   │   │   ├── create.gorm.go           🔲
│   │   │   ├── update.gorm.go           🔲
│   │   │   ├── delete.gorm.go           🔲 — soft delete
│   │   │   ├── find_by_id.gorm.go       🔲
│   │   │   └── find_all.gorm.go         🔲 — with filters, pagination
│   │   ├── genre/                       🔲 TODO
│   │   │   ├── genre.gorm.go            🔲
│   │   │   ├── create.gorm.go           🔲
│   │   │   ├── delete.gorm.go           🔲
│   │   │   └── find_all.gorm.go         🔲
│   │   ├── theater/                     🔲 TODO
│   │   │   ├── theater.gorm.go          🔲
│   │   │   ├── create.gorm.go           🔲
│   │   │   ├── update.gorm.go           🔲
│   │   │   ├── find_by_id.gorm.go       🔲
│   │   │   └── find_all.gorm.go         🔲
│   │   ├── seat/                        🔲 TODO
│   │   │   ├── seat.gorm.go             🔲
│   │   │   ├── create_batch.gorm.go     🔲 — auto-generate on theater creation
│   │   │   └── find_by_theater_id.gorm.go 🔲
│   │   ├── showtime/                    🔲 TODO
│   │   │   ├── showtime.gorm.go         🔲
│   │   │   ├── create.gorm.go           🔲
│   │   │   ├── update.gorm.go           🔲
│   │   │   ├── delete.gorm.go           🔲
│   │   │   ├── find_by_id.gorm.go       🔲
│   │   │   ├── find_by_date.gorm.go     🔲
│   │   │   └── check_conflict.gorm.go   🔲 — overlap detection
│   │   ├── reservation/                 🔲 TODO
│   │   │   ├── reservation.gorm.go      🔲
│   │   │   ├── create.gorm.go           🔲 — within TX: lock seats + insert
│   │   │   ├── cancel.gorm.go           🔲
│   │   │   ├── find_by_id.gorm.go       🔲
│   │   │   ├── find_by_user_id.gorm.go  🔲
│   │   │   └── find_all.gorm.go         🔲 — admin: all reservations with filters
│   │   └── reservation_seat/            🔲 TODO
│   │       ├── reservation_seat.gorm.go 🔲
│   │       ├── create_batch.gorm.go     🔲
│   │       ├── find_by_showtime_id.gorm.go 🔲 — seat availability map
│   │       └── check_availability.gorm.go  🔲 — SELECT FOR UPDATE lock check
│   └── mock/                            ← Auto-generated gomock mocks (make repository-mocks)
│       ├── user.repository_mock.go      ✅
│       ├── refresh_token.repository_mock.go ✅
│       ├── counter.repository_mock.go   ✅
│       ├── message.repository_mock.go   ✅
│       ├── movie.repository_mock.go     🔲 TODO — auto-generated
│       ├── genre.repository_mock.go     🔲 TODO — auto-generated
│       ├── theater.repository_mock.go   🔲 TODO — auto-generated
│       ├── seat.repository_mock.go      🔲 TODO — auto-generated
│       ├── showtime.repository_mock.go  🔲 TODO — auto-generated
│       ├── reservation.repository_mock.go 🔲 TODO — auto-generated
│       └── reservation_seat.repository_mock.go 🔲 TODO — auto-generated
│
└── service/                             ← Business logic (1 package per domain, 1 file per method)
    ├── token/                           ✅ Exists — standalone (no repo dependency)
    │   ├── token.service.go             ✅ — TokenService interface + TokenConfig
    │   ├── generate_access.service.go   ✅
    │   ├── generate_refresh.service.go  ✅
    │   ├── validate_access.service.go   ✅
    │   ├── validate_refresh.service.go  ✅
    │   └── *_test.go                    ✅
    ├── user/                            ✅ Exists
    │   ├── user.service.go              ✅ — UserService interface
    │   ├── register.service.go          ✅
    │   ├── login.service.go             ✅
    │   ├── refresh.service.go           ✅
    │   ├── store_refresh_token.service.go ✅
    │   ├── revoke_refresh_tokens.service.go ✅
    │   └── *_test.go                    ✅
    ├── csrf/                            ✅ Exists
    │   ├── csrf.service.go              ✅
    │   ├── generate.service.go          ✅
    │   ├── validate.service.go          ✅
    │   └── *_test.go                    ✅
    ├── health/                          ✅ Exists
    ├── counter/                         ✅ Exists — demo
    ├── message/                         ✅ Exists — demo
    ├── movie/                           🔲 TODO
    │   ├── movie.service.go             🔲 — MovieService interface + constructor
    │   ├── create_movie.service.go      🔲 — validate + create + attach genres
    │   ├── update_movie.service.go      🔲 — validate + update + reassign genres
    │   ├── delete_movie.service.go      🔲 — soft delete (check no future active reservations)
    │   ├── get_movie.service.go         🔲 — find by ID with genres
    │   ├── list_movies.service.go       🔲 — filter by genre/date/title + pagination
    │   └── *_test.go                    🔲
    ├── genre/                           🔲 TODO
    │   ├── genre.service.go             🔲
    │   ├── create_genre.service.go      🔲
    │   ├── delete_genre.service.go      🔲 — check no movies tagged
    │   ├── list_genres.service.go       🔲
    │   └── *_test.go                    🔲
    ├── theater/                         🔲 TODO
    │   ├── theater.service.go           🔲
    │   ├── create_theater.service.go    🔲 — auto-generate seats (A1..Z30)
    │   ├── update_theater.service.go    🔲 — name/location only
    │   ├── get_theater.service.go       🔲 — with seat layout
    │   ├── list_theaters.service.go     🔲
    │   └── *_test.go                    🔲
    ├── showtime/                        🔲 TODO
    │   ├── showtime.service.go          🔲
    │   ├── create_showtime.service.go   🔲 — conflict detection + auto end_time
    │   ├── update_showtime.service.go   🔲 — check no confirmed reservations
    │   ├── delete_showtime.service.go   🔲 — auto-cancel reservations if confirmed
    │   ├── get_showtime.service.go      🔲
    │   ├── list_showtimes.service.go    🔲 — by date, grouped by movie
    │   ├── get_seat_map.service.go      🔲 — available/reserved status per seat
    │   └── *_test.go                    🔲
    ├── reservation/                     🔲 TODO
    │   ├── reservation.service.go       🔲
    │   ├── create_reservation.service.go 🔲 — TX: lock seats → check → insert → booking_reference
    │   ├── cancel_reservation.service.go 🔲 — ownership check + 1hr window
    │   ├── get_reservation.service.go   🔲
    │   ├── list_my_reservations.service.go 🔲 — by user, filter by status
    │   └── *_test.go                    🔲
    └── admin/                           🔲 TODO
        ├── admin.service.go             🔲
        ├── list_all_reservations.service.go 🔲 — filters: user, movie, showtime, status, date
        ├── cancel_any_reservation.service.go 🔲 — no timing restriction
        ├── promote_user.service.go      🔲 — prevent self-demotion
        ├── revenue_report.service.go    🔲 — grouped by movie/theater/date
        ├── capacity_report.service.go   🔲 — occupancy per showtime
        ├── dashboard_summary.service.go 🔲 — aggregated stats
        └── *_test.go                    🔲
```

### Frontend

| Layer              | Technology                     |
|--------------------|--------------------------------|
| Framework          | Next.js 14 (App Router)        |
| Language           | TypeScript                     |
| Styling            | Tailwind CSS                   |
| UI Components      | shadcn/ui                      |
| State Management   | Zustand + React Query          |
| HTTP Client        | Axios                          |
| Forms              | React Hook Form + Zod          |
| Date Handling      | date-fns                       |

### Infrastructure

```
┌────────────────┐     HTTP/REST     ┌──────────────────┐
│  Next.js       │ ───────────────► │  Go API Server   │
│  (Frontend)    │                  │  (Gin)           │
│  Port 3000     │                  │  Port 8080       │
└────────────────┘                  └────────┬─────────┘
                                             │ GORM
                                    ┌────────▼─────────┐
                                    │ PostgreSQL / SQLite│
                                    │ (configurable)    │
                                    └──────────────────┘
```

**Development:** SQLite (zero-config, `DATABASE_TYPE=sqlite`, default `dev.db`)
**Production:** PostgreSQL (`DATABASE_TYPE=postgres`, `DATABASE_DSN=<connection_string>`)

Entry point: `cmd/server/main.go` → loads `.env` via godotenv → builds DI container → starts HTTP server with graceful shutdown.

Hot-reload in development via Air (`make dev`).

---

## 11. Security Considerations

| Threat                  | Mitigation                                                              |
|-------------------------|-------------------------------------------------------------------------|
| Brute force login       | Rate limiting on `/auth/login` (10 req/min per IP)                     |
| JWT token theft         | Short expiry (24h), HTTPS only, HttpOnly cookies optional               |
| SQL injection           | GORM parameterized queries; no raw string interpolation                 |
| IDOR (Insecure Direct)  | Ownership checks: `reservation.user_id == jwt.user_id`                 |
| Password exposure       | bcrypt hashed; never returned in responses                              |
| Mass assignment         | Explicit DTO structs; GORM `Select` / `Omit` to prevent over-posting    |
| CORS misconfiguration   | Whitelist specific frontend origins; credentials: true only if needed   |
| Overbooking race        | DB transaction + `SELECT FOR UPDATE` + unique partial index             |
| Admin privilege escalation | Only admins can promote users; self-demotion prevented              |

---

## 12. Reporting & Analytics

### Revenue Report Response
```json
{
  "summary": {
    "total_revenue": 45230000,
    "total_tickets": 1845,
    "period": "2025-05-01 to 2025-05-31"
  },
  "by_movie": [
    {
      "movie_id": 3,
      "movie_title": "Interstellar",
      "tickets_sold": 420,
      "revenue": 12600000
    }
  ]
}
```

### Capacity Report Response
```json
{
  "showtimes": [
    {
      "showtime_id": 101,
      "movie_title": "Dune: Part Two",
      "theater_name": "Hall B",
      "start_time": "2025-06-01T14:00:00Z",
      "total_seats": 80,
      "seats_reserved": 72,
      "seats_available": 8,
      "occupancy_rate": 90.0
    }
  ]
}
```

---

## 13. Out of Scope

The following are **explicitly excluded** from v1.0 to keep scope manageable:

- Payment processing (Stripe, Midtrans, etc.)
- Email/SMS notifications
- Waiting list / queue for sold-out showtimes
- Dynamic seat pricing (surge pricing)
- Mobile application (iOS/Android)
- Multi-language (i18n) support
- Social login (Google/Facebook OAuth)
- User ratings and reviews for movies
- Print/PDF ticket generation

---

## 14. Future Enhancements

| Priority | Feature                           | Notes                                          |
|----------|-----------------------------------|------------------------------------------------|
| High     | Payment integration               | Midtrans or Stripe for IDR/USD                 |
| High     | Email confirmation                | SendGrid or AWS SES on reservation confirmed   |
| Medium   | Seat hold / temporary lock        | 10-minute lock during seat selection flow      |
| Medium   | User reviews & ratings            | Stars + text review per movie                  |
| Medium   | QR Code ticket                    | Scannable QR for entry validation              |
| Low      | Loyalty points                    | Earn points per reservation, redeem discounts  |
| Low      | Recommendation engine             | "You might also like" based on history         |
| Low      | Real-time seat updates            | WebSocket to reflect seat status live          |

---

## 15. Glossary

| Term              | Definition                                                           |
|-------------------|----------------------------------------------------------------------|
| **Showtime**      | A scheduled screening of a movie at a specific theater and time     |
| **Seat Map**      | Visual grid representation of all seats in a theater                |
| **Booking Ref**   | Human-readable unique reservation identifier (e.g., CINE-20250601-00042) |
| **Overbooking**   | State where more reservations exist than available seats (must never occur) |
| **Soft Delete**   | Marking a record as deleted without physically removing from DB     |
| **RBAC**          | Role-Based Access Control — permissions tied to user roles          |
| **JWT**           | JSON Web Token — stateless authentication token                     |
| **DTO**           | Data Transfer Object — struct used for API request/response shaping |
| **FOR UPDATE**    | PostgreSQL clause to lock rows during a transaction                 |
| **Occupancy Rate**| `(reserved_seats / total_seats) × 100%`                            |

---

*End of Document — Movie Reservation System PRD v1.0*