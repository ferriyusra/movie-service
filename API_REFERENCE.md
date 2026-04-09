# API Reference — Movie Reservation System

Base URL: `http://localhost:8080/api`

## Authentication

All protected endpoints require the `Authorization` header:

```
Authorization: Bearer <accessToken>
```

**Roles:**
- `user` — can browse, reserve seats, manage own reservations
- `admin` — full access including CRUD management and reports

**Seeded admin account:**
```
email:    admin@cinema.com
password: Admin@123
```

---

## Response Envelope

Every response follows this structure:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 150
  },
  "errors": {
    "email": "Email is required"
  }
}
```

| Field     | Type              | Description                                   |
|-----------|-------------------|-----------------------------------------------|
| `success` | `boolean`         | `true` if the request succeeded               |
| `message` | `string`          | Human-readable status message                 |
| `data`    | `object \| array` | Response payload (omitted on error)           |
| `meta`    | `object \| null`  | Pagination info (only on paginated endpoints) |
| `errors`  | `object \| null`  | Field-level validation errors                 |

---

## TypeScript Interfaces

Copy these into your frontend project for type safety:

```typescript
// ─── Response Envelope ───────────────────────────────────────
interface ApiResponse<T> {
  success: boolean;
  message: string;
  data?: T;
  meta?: Meta;
  errors?: Record<string, string>;
}

interface Meta {
  page: number;
  limit: number;
  total: number;
}

// ─── Auth ────────────────────────────────────────────────────
interface User {
  id: string;
  email: string;
  name: string;
}

interface LoginResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
}

interface RegisterResponse {
  user: User;
}

interface RefreshResponse {
  accessToken: string;
  refreshToken: string;
}

// ─── Genre ───────────────────────────────────────────────────
interface Genre {
  id: string;
  name: string;
}

// ─── Movie ───────────────────────────────────────────────────
interface Movie {
  id: string;
  title: string;
  description: string;
  posterUrl: string;
  durationMin: number;
  language: string;
  releaseDate: string | null;
  rating: number;
  genres: Genre[];
  createdAt: string;
  updatedAt: string;
}

// ─── Theater & Seat ──────────────────────────────────────────
interface Theater {
  id: string;
  name: string;
  location: string;
  totalSeats: number;
  rows: number;
  seatsPerRow: number;
}

interface Seat {
  id: string;
  row: string;
  number: number;
  label: string;
  type: "standard" | "vip" | "accessible";
}

interface TheaterDetail extends Theater {
  seats: Seat[];
}

// ─── Showtime ────────────────────────────────────────────────
interface Showtime {
  id: string;
  movieId: string;
  theaterId: string;
  startTime: string;
  endTime: string;
  price: number;
  availableSeats: number;
}

interface ShowtimeDetail extends Showtime {
  movie: Movie;
  theater: Theater;
}

// ─── Reservation ─────────────────────────────────────────────
interface Reservation {
  id: string;
  bookingReference: string;
  showtimeId: string;
  status: "confirmed" | "cancelled";
  totalAmount: number;
  cancelledAt?: string;
  createdAt: string;
  seats: Seat[];
}

interface ReservationDetail extends Reservation {
  movie: Movie;
  showtime: Showtime;
  theater: Theater;
}
```

---

## Endpoints

### Auth

#### `POST /api/auth/register`
Public — Create a new user account.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "Password123"
}
```

**Response** `201`:
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "user": {
      "id": "uuid",
      "email": "john@example.com",
      "name": "John Doe"
    }
  }
}
```

---

#### `POST /api/auth/login`
Public — Authenticate and receive tokens.

**Request Body:**
```json
{
  "email": "admin@cinema.com",
  "password": "Admin@123"
}
```

**Response** `200`:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": { "id": "uuid", "email": "admin@cinema.com", "name": "Admin" },
    "accessToken": "eyJhbG...",
    "refreshToken": "eyJhbG..."
  }
}
```

---

#### `POST /api/auth/refresh`
Public — Exchange refresh token for a new token pair.

**Request Body:**
```json
{
  "refreshToken": "eyJhbG..."
}
```

**Response** `200`:
```json
{
  "success": true,
  "data": {
    "accessToken": "eyJhbG...",
    "refreshToken": "eyJhbG..."
  }
}
```

---

#### `GET /api/auth/me`
Protected — Get the current user's profile.

**Response** `200`:
```json
{
  "success": true,
  "data": { "id": "uuid", "email": "john@example.com", "name": "John Doe" }
}
```

---

#### `POST /api/auth/logout`
Protected — Revoke refresh tokens.

**Response** `200`:
```json
{ "success": true, "message": "Logged out successfully" }
```

---

### Genres

#### `GET /api/genres`
Public — List all genres.

**Response** `200`:
```json
{
  "success": true,
  "data": [
    { "id": "uuid", "name": "Action" },
    { "id": "uuid", "name": "Comedy" }
  ]
}
```

---

#### `POST /api/genres`
Admin — Create a genre.

**Request Body:**
```json
{ "name": "Sci-Fi" }
```

**Response** `201`:
```json
{ "success": true, "data": { "id": "uuid", "name": "Sci-Fi" } }
```

---

#### `DELETE /api/genres/:id`
Admin — Delete a genre (fails if movies are tagged with it).

**Response** `200`:
```json
{ "success": true, "message": "Genre deleted" }
```

---

### Movies

#### `GET /api/movies`
Public — Browse movies with optional filters and pagination.

**Query Params:**
| Param     | Type   | Default | Description              |
|-----------|--------|---------|--------------------------|
| `page`    | number | `1`     | Page number              |
| `limit`   | number | `20`    | Items per page (max 100) |
| `title`   | string | —       | Search by title          |
| `genreId` | UUID   | —       | Filter by genre          |

**Response** `200`:
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "title": "Inception",
      "description": "A mind-bending thriller",
      "posterUrl": "https://example.com/poster.jpg",
      "durationMin": 148,
      "language": "English",
      "releaseDate": "2010-07-16T00:00:00Z",
      "rating": 8.8,
      "genres": [{ "id": "uuid", "name": "Sci-Fi" }],
      "createdAt": "2026-04-09T10:00:00Z",
      "updatedAt": "2026-04-09T10:00:00Z"
    }
  ],
  "meta": { "page": 1, "limit": 20, "total": 42 }
}
```

---

#### `GET /api/movies/:id`
Public — Get a single movie with genres.

**Response** `200`: Single `Movie` object in `data`.

---

#### `POST /api/movies`
Admin — Create a movie.

**Request Body:**
```json
{
  "title": "Inception",
  "description": "A mind-bending thriller",
  "posterUrl": "https://example.com/poster.jpg",
  "durationMin": 148,
  "language": "English",
  "releaseDate": "2010-07-16",
  "genreIds": ["<genre-uuid-1>", "<genre-uuid-2>"]
}
```

**Response** `201`: `Movie` object in `data`.

---

#### `PUT /api/movies/:id`
Admin — Update a movie (same body as create).

**Response** `200`: Updated `Movie` object in `data`.

---

#### `DELETE /api/movies/:id`
Admin — Soft-delete a movie.

**Response** `200`.

---

### Theaters

#### `GET /api/theaters`
Public — List all theaters.

**Response** `200`:
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Theater 1",
      "location": "Ground Floor",
      "totalSeats": 50,
      "rows": 5,
      "seatsPerRow": 10
    }
  ]
}
```

---

#### `GET /api/theaters/:id`
Public — Get theater detail with full seat layout.

**Response** `200`:
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Theater 1",
    "location": "Ground Floor",
    "totalSeats": 50,
    "rows": 5,
    "seatsPerRow": 10,
    "seats": [
      { "id": "uuid", "row": "A", "number": 1, "label": "A1", "type": "standard" },
      { "id": "uuid", "row": "A", "number": 2, "label": "A2", "type": "standard" }
    ]
  }
}
```

---

#### `POST /api/theaters`
Admin — Create a theater (seats are auto-generated).

**Request Body:**
```json
{
  "name": "Theater 1",
  "location": "Ground Floor",
  "rows": 5,
  "seatsPerRow": 10
}
```

**Response** `201`: `TheaterDetail` object with generated `seats` array.

---

#### `PUT /api/theaters/:id`
Admin — Update theater name and location only.

**Request Body:**
```json
{
  "name": "Theater 1 (Renovated)",
  "location": "Ground Floor"
}
```

**Response** `200`: `Theater` object.

---

### Showtimes

#### `GET /api/showtimes`
Public — List showtimes for a date.

**Query Params:**
| Param  | Type   | Default | Description                   |
|--------|--------|---------|-------------------------------|
| `date` | string | today   | Date filter (`YYYY-MM-DD`)    |

**Response** `200`: Array of `Showtime` objects.

---

#### `GET /api/showtimes/:id`
Public — Get showtime detail with movie and theater info.

**Response** `200`: `ShowtimeDetail` object (includes nested `movie` and `theater`).

---

#### `POST /api/showtimes`
Admin — Create a showtime. `endTime` is auto-calculated as `startTime + durationMin + 15min`.

**Request Body:**
```json
{
  "movieId": "<movie-uuid>",
  "theaterId": "<theater-uuid>",
  "startTime": "2026-04-10T14:00:00Z",
  "price": 50000
}
```

**Response** `201`: `Showtime` object.

**Error** `409` — Conflict if the time slot overlaps with an existing showtime in the same theater.

---

#### `PUT /api/showtimes/:id`
Admin — Update a showtime (same body as create).

**Response** `200`: `Showtime` object.

---

#### `DELETE /api/showtimes/:id`
Admin — Soft-delete a showtime.

**Response** `200`.

---

### Reservations

#### `POST /api/reservations`
Protected — Reserve seats for a showtime.

**Request Body:**
```json
{
  "showtimeId": "<showtime-uuid>",
  "seatIds": ["<seat-uuid-1>", "<seat-uuid-2>"]
}
```

**Constraints:**
- 1–10 seats per reservation
- Cannot reserve past showtimes
- Seats must be available (not already reserved for this showtime)

**Response** `201`:
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "bookingReference": "CINE-20260410-a1b2c",
    "showtimeId": "uuid",
    "status": "confirmed",
    "totalAmount": 100000,
    "createdAt": "2026-04-09T12:00:00Z",
    "seats": [
      { "id": "uuid", "row": "A", "number": 1, "label": "A1", "type": "standard" },
      { "id": "uuid", "row": "A", "number": 2, "label": "A2", "type": "standard" }
    ]
  }
}
```

**Error** `409` — Seats already reserved.

---

#### `GET /api/reservations`
Protected — List the current user's reservations.

**Query Params:**
| Param    | Type   | Default | Description                          |
|----------|--------|---------|--------------------------------------|
| `page`   | number | `1`     | Page number                          |
| `limit`  | number | `20`    | Items per page (max 100)             |
| `status` | string | —       | Filter: `confirmed` or `cancelled`   |

**Response** `200`: Array of `Reservation` objects with `meta`.

---

#### `GET /api/reservations/:id`
Protected — Get reservation detail (ownership enforced).

**Response** `200`: `ReservationDetail` object (includes `movie`, `showtime`, `theater`).

**Error** `403` — Not your reservation.

---

#### `DELETE /api/reservations/:id`
Protected — Cancel a reservation.

**Constraints:**
- Must be the owner
- Must be at least 1 hour before showtime `startTime`
- Must not already be cancelled

**Response** `200`.

**Error** `403` — Not your reservation.
**Error** `422` — Cancellation window has passed.

---

### Admin

> All admin endpoints require `Authorization` header with an admin token.

#### `GET /api/admin/reservations`
List all reservations with filters.

**Query Params:**
| Param        | Type   | Default | Description               |
|--------------|--------|---------|---------------------------|
| `page`       | number | `1`     | Page number               |
| `limit`      | number | `20`    | Items per page (max 100)  |
| `status`     | string | —       | `confirmed` / `cancelled` |
| `userId`     | UUID   | —       | Filter by user            |
| `showtimeId` | UUID   | —       | Filter by showtime        |
| `dateFrom`   | string | —       | Created after (date)      |
| `dateTo`     | string | —       | Created before (date)     |

**Response** `200`: Array of `Reservation` objects with `meta`.

---

#### `DELETE /api/admin/reservations/:id`
Force-cancel any reservation (no time restriction).

**Response** `200`.

---

#### `GET /api/admin/users`
List all users.

**Response** `200`:
```json
{
  "success": true,
  "data": [
    { "id": "uuid", "email": "admin@cinema.com", "name": "Admin" },
    { "id": "uuid", "email": "john@example.com", "name": "John Doe" }
  ]
}
```

---

#### `PATCH /api/admin/users/:id/promote`
Promote a user to admin role. Cannot promote yourself.

**Response** `200`.

**Error** `400` — "cannot promote yourself" or "user is already an admin".

---

## HTTP Status Codes

| Code  | Meaning                                                  |
|-------|----------------------------------------------------------|
| `200` | Success                                                  |
| `201` | Created (POST requests that create a resource)           |
| `400` | Bad request / validation error                           |
| `401` | Missing or invalid authentication token                  |
| `403` | Forbidden (wrong role or not the resource owner)         |
| `404` | Resource not found                                       |
| `409` | Conflict (seat already reserved, showtime overlap)       |
| `422` | Unprocessable (cancellation window passed)               |
| `500` | Internal server error                                    |

---

## Frontend Integration Tips

### Axios Instance Setup

```typescript
import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080/api",
  headers: { "Content-Type": "application/json" },
});

// Attach access token to every request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("accessToken");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auto-refresh on 401
api.interceptors.response.use(
  (res) => res,
  async (error) => {
    const originalRequest = error.config;
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      const refreshToken = localStorage.getItem("refreshToken");
      if (refreshToken) {
        try {
          const { data } = await axios.post(
            "http://localhost:8080/api/auth/refresh",
            { refreshToken }
          );
          localStorage.setItem("accessToken", data.data.accessToken);
          localStorage.setItem("refreshToken", data.data.refreshToken);
          originalRequest.headers.Authorization = `Bearer ${data.data.accessToken}`;
          return api(originalRequest);
        } catch {
          localStorage.clear();
          window.location.href = "/login";
        }
      }
    }
    return Promise.reject(error);
  }
);

export default api;
```

### Typical API Calls

```typescript
// Login
const { data } = await api.post<ApiResponse<LoginResponse>>("/auth/login", {
  email: "admin@cinema.com",
  password: "Admin@123",
});
localStorage.setItem("accessToken", data.data!.accessToken);
localStorage.setItem("refreshToken", data.data!.refreshToken);

// List movies with pagination
const { data } = await api.get<ApiResponse<Movie[]>>("/movies", {
  params: { page: 1, limit: 20, genreId: "some-uuid" },
});
const movies = data.data!;
const total = data.meta!.total;

// Create a reservation
const { data } = await api.post<ApiResponse<Reservation>>("/reservations", {
  showtimeId: "some-uuid",
  seatIds: ["seat-uuid-1", "seat-uuid-2"],
});
const bookingRef = data.data!.bookingReference;

// Cancel a reservation
await api.delete(`/reservations/${reservationId}`);
```

### Seat Picker Flow

1. **Get theater seats:** `GET /api/theaters/:id` → use `seats` array to render the grid
2. **Get reserved seats:** Fetch showtime detail or use seat availability from the reservation flow
3. **Submit reservation:** `POST /api/reservations` with selected `seatIds`
4. **Handle 409:** Show which seats were taken and refresh the seat map

### Auth Guard (Next.js Middleware)

```typescript
// middleware.ts
import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  const token = request.cookies.get("accessToken")?.value;
  const isAuthPage = request.nextUrl.pathname.startsWith("/login") ||
                     request.nextUrl.pathname.startsWith("/register");
  const isAdminPage = request.nextUrl.pathname.startsWith("/admin");

  if (!token && !isAuthPage) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  // For admin pages, decode JWT and check role on the client side
  // or use a server component to verify
  return NextResponse.next();
}

export const config = {
  matcher: ["/reservations/:path*", "/admin/:path*", "/profile/:path*"],
};
```
