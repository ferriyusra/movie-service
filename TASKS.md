# Task List — Movie Reservation System

Berdasarkan [PRD.md](PRD.md) v1.0. Setiap fase backend harus selesai sebelum frontend terkait bisa dibangun.

---

## Fase 1: Foundation & Core Models

### Backend

- [ ] **1. Entity models** — Buat semua GORM entities: `movie`, `genre`, `movie_genre`, `theater`, `seat`, `showtime`, `reservation`, `reservation_seat` di `internal/model/entity/`
- [ ] **2. Request DTOs** — Buat request structs di `internal/model/request/`: `movie.go`, `genre.go`, `theater.go`, `showtime.go`, `reservation.go`
- [ ] **3. Response DTOs** — Buat response structs di `internal/model/response/`: `movie.go`, `genre.go`, `theater.go`, `showtime.go`, `reservation.go`, `report.go`
- [ ] **4. Database auto-migrate** — Tambahkan semua entity baru ke auto-migrate di `internal/platform/database.go`
- [ ] **5. RBAC middleware** — Buat `AdminOnly` middleware di `internal/api/middleware/rbac.go` (cek `role = admin` dari JWT claims)
- [ ] **6. Seed admin account** — Bootstrap `admin@cinema.com / Admin@123` dengan role `admin` saat startup pertama

### Frontend

- [ ] **7. Project setup** — Init Next.js 14 (App Router) + TypeScript + Tailwind CSS + shadcn/ui
- [ ] **8. Axios instance** — Setup HTTP client dengan interceptor untuk JWT access/refresh token
- [ ] **9. Auth store (Zustand)** — State management untuk user session (login, logout, user info)
- [ ] **10. Layout** — Public layout (navbar, footer) + Admin layout (sidebar navigation)

---

## Fase 2: Genre & Movie Management

### Backend

- [ ] **11. Genre repository interface** — `internal/repository/interfaces/genre.repository_interface.go` — Create, Delete, FindAll, FindByID
- [ ] **12. Genre GORM implementation** — `internal/repository/implementations/genre/` — genre.gorm.go, create.gorm.go, delete.gorm.go, find_all.gorm.go
- [ ] **13. Genre service** — `internal/service/genre/` — CreateGenre, DeleteGenre (cek no movies tagged), ListGenres
- [ ] **14. Genre handler + routes** — `POST /genres` (admin), `DELETE /genres/:id` (admin), `GET /genres` (public)
- [ ] **15. Genre tests** — Service unit tests dengan gomock (table-driven)
- [ ] **16. Movie repository interface** — `internal/repository/interfaces/movie.repository_interface.go` — Create, Update, Delete, FindByID, FindAll (filters + pagination)
- [ ] **17. Movie GORM implementation** — `internal/repository/implementations/movie/` — semua method files
- [ ] **18. Movie service** — `internal/service/movie/` — CreateMovie (+ attach genres), UpdateMovie (+ reassign genres), DeleteMovie (soft, cek future reservations), GetMovie, ListMovies
- [ ] **19. Movie handler + routes** — CRUD admin (`POST/PUT/DELETE /movies`) + public (`GET /movies`, `GET /movies/:id`)
- [ ] **20. Movie tests** — Service unit tests
- [ ] **21. Generate mocks** — `make repository-mocks` untuk genre & movie interfaces
- [ ] **22. DI wiring** — Tambahkan genre & movie ke `internal/di/container.go`

### Frontend

- [ ] **23. Genre list page** — `/admin/genres` — CRUD table (create, delete)
- [ ] **24. Movie list page (public)** — `/movies` — Browsing dengan genre filter + search + pagination
- [ ] **25. Movie detail page** — `/movies/[id]` — Synopsis, genres, poster, duration, language
- [ ] **26. Admin movie list** — `/admin/movies` — Table dengan edit/delete actions
- [ ] **27. Create movie form** — `/admin/movies/new` — Form dengan React Hook Form + Zod validation
- [ ] **28. Edit movie form** — `/admin/movies/[id]/edit` — Pre-filled form

---

## Fase 3: Theater & Seat Management

### Backend

- [ ] **29. Theater repository interface** — Create, Update, FindByID, FindAll
- [ ] **30. Seat repository interface** — CreateBatch, FindByTheaterID
- [ ] **31. Theater GORM implementation** — `internal/repository/implementations/theater/`
- [ ] **32. Seat GORM implementation** — `internal/repository/implementations/seat/` — auto-generate seats (A1..Z30)
- [ ] **33. Theater service** — CreateTheater (auto-generate seats), UpdateTheater (name/location only, cek no past reservations untuk layout change), GetTheater (with seat layout), ListTheaters
- [ ] **34. Theater handler + routes** — CRUD admin + public list/detail
- [ ] **35. Theater tests** — Service unit tests
- [ ] **36. Generate mocks** — `make repository-mocks` untuk theater & seat
- [ ] **37. DI wiring** — Tambahkan theater & seat ke container

### Frontend

- [ ] **38. Admin theater management** — `/admin/theaters` — Create/edit theater, lihat seat layout grid
- [ ] **39. Theater detail (public)** — `/theaters/:id` — View theater info + seat grid visualization

---

## Fase 4: Showtime Management

### Backend

- [ ] **40. Showtime repository interface** — Create, Update, Delete, FindByID, FindByDate, CheckConflict
- [ ] **41. Showtime GORM implementation** — `internal/repository/implementations/showtime/` — termasuk `check_conflict.gorm.go` (overlap detection)
- [ ] **42. Showtime service** — CreateShowtime (conflict detection + auto `end_time = start_time + duration + 15min`), UpdateShowtime (cek no confirmed reservations), DeleteShowtime (auto-cancel reservations), GetShowtime, ListShowtimes by date (grouped by movie), GetSeatMap (available/reserved per seat)
- [ ] **43. Showtime handler + routes** — CRUD admin + public (`GET /showtimes?date=`, `GET /showtimes/:id`, `GET /showtimes/:id/seats`)
- [ ] **44. Showtime tests** — Service unit tests (terutama conflict detection edge cases)
- [ ] **45. Generate mocks** — `make repository-mocks` untuk showtime
- [ ] **46. DI wiring** — Tambahkan showtime ke container

### Frontend

- [ ] **47. Showtime list di movie detail** — Tampilkan showtimes by date di halaman `/movies/[id]`
- [ ] **48. Seat picker page** — `/showtimes/[id]/seats` — Interactive seat map (green=available, red=reserved, blue=selected), max 10 seats
- [ ] **49. Admin showtime CRUD** — `/admin/showtimes` — Calendar/list view + create/edit/delete forms

---

## Fase 5: Reservation System (Core Feature)

### Backend

- [ ] **50. Reservation repository interface** — Create (within TX), Cancel, FindByID, FindByUserID, FindAll
- [ ] **51. ReservationSeat repository interface** — CreateBatch, FindByShowtimeID, CheckAvailability (`SELECT FOR UPDATE`)
- [ ] **52. Reservation GORM implementation** — `internal/repository/implementations/reservation/` — `create.gorm.go` harus dalam single DB transaction (lock seats -> check -> insert)
- [ ] **53. ReservationSeat GORM implementation** — `internal/repository/implementations/reservation_seat/`
- [ ] **54. Reservation service** — CreateReservation (TX: lock seats -> check availability -> insert reservation + seats -> generate booking_reference `CINE-YYYYMMDD-XXXXX` -> calculate total_amount), CancelReservation (ownership check + 1hr window before showtime), GetReservation, ListMyReservations (filter by status)
- [ ] **55. Reservation handler + routes** — `POST /reservations` (user), `DELETE /reservations/:id` (user), `GET /reservations` (user), `GET /reservations/:id` (user)
- [ ] **56. Reservation tests** — Service unit tests (terutama concurrency/overbooking scenarios)
- [ ] **57. DB partial unique index** — `CREATE UNIQUE INDEX idx_no_double_booking` untuk prevent double-booking di level database
- [ ] **58. Generate mocks** — `make repository-mocks` untuk reservation & reservation_seat
- [ ] **59. DI wiring** — Tambahkan reservation ke container

### Frontend

- [ ] **60. Booking confirmation flow** — Dari seat picker -> review selected seats -> confirm -> tampilkan booking reference + summary
- [ ] **61. My reservations page** — `/reservations` — List dengan filter status (confirmed/cancelled) + pagination
- [ ] **62. Reservation detail page** — `/reservations/[id]` — Movie, showtime, theater, seats, total amount, status + cancel button
- [ ] **63. Cancel reservation flow** — Dialog konfirmasi + handle error responses (422 window passed, 403 not owner)

---

## Fase 6: Admin Features & Reporting

### Backend

- [ ] **64. Admin service** — ListAllReservations (filters: user_id, movie_id, showtime_id, status, date_from, date_to), CancelAnyReservation (no time restriction), PromoteUser (prevent self-demotion)
- [ ] **65. Revenue report service** — Total revenue grouped by movie/theater/date range. Response: `movie_title, showtime_count, tickets_sold, total_revenue`
- [ ] **66. Capacity report service** — Per showtime: `total_seats, seats_reserved, seats_available, occupancy_rate (%)`
- [ ] **67. Dashboard summary service** — Total movies/showtimes/users/active reservations today, revenue this week vs last week, top 5 movies by tickets sold, top 5 by revenue
- [ ] **68. Admin handler + routes** — `GET /admin/reservations`, `DELETE /admin/reservations/:id`, `GET /admin/reports/revenue`, `GET /admin/reports/capacity`, `GET /admin/reports/dashboard`, `GET /admin/users`, `PATCH /admin/users/:id/promote`
- [ ] **69. Admin tests** — Service unit tests
- [ ] **70. Rate limiting** — 10 req/min per IP pada `POST /auth/login`

### Frontend

- [ ] **71. Admin dashboard** — `/admin` — Stats cards (total movies, users, reservations, revenue), revenue chart (week comparison), top 5 movies tables
- [ ] **72. Admin reservations** — `/admin/reservations` — Data table dengan filters (user, movie, showtime, status, date range) + force cancel
- [ ] **73. Admin reports page** — `/admin/reports` — Revenue report table/chart + capacity report table (occupancy bars)
- [ ] **74. Admin users page** — `/admin/users` — User list table + promote to admin button with confirmation
- [ ] **75. User profile page** — `/profile` — View/edit nama, email

---

## Fase 7: Polish & Security

### Backend

- [ ] **76. Structured JSON logging** — Request logger middleware: method, path, status code, latency, user_id (if authenticated)
- [ ] **77. Input validation** — go-playground/validator tags untuk semua request DTOs (min length, required, email format, etc.)
- [ ] **78. Pagination guard** — Enforce max `limit=100` dan default `limit=20` untuk semua list endpoints
- [ ] **79. CORS config** — Configurable `ALLOWED_ORIGINS` dari env, credentials support

### Frontend

- [ ] **80. Error handling** — Global error boundary + toast notifications (success/error) untuk semua API calls
- [ ] **81. Loading states** — Skeleton loaders / spinners untuk semua data fetching pages
- [ ] **82. Responsive design** — Mobile-friendly layout untuk semua pages (terutama seat picker)
- [ ] **83. Auth guard** — Redirect ke `/login` untuk protected routes, redirect ke `/` jika bukan admin akses `/admin/*`
- [ ] **84. Login & Register pages** — `/login` + `/register` — Forms dengan validation (React Hook Form + Zod)

---

## Summary

| Fase | Backend Tasks | Frontend Tasks | Total |
|------|:------------:|:--------------:|:-----:|
| 1. Foundation & Core Models | 6 | 4 | 10 |
| 2. Genre & Movie Management | 12 | 6 | 18 |
| 3. Theater & Seat Management | 9 | 2 | 11 |
| 4. Showtime Management | 7 | 3 | 10 |
| 5. Reservation System | 10 | 4 | 14 |
| 6. Admin & Reporting | 7 | 5 | 12 |
| 7. Polish & Security | 4 | 5 | 9 |
| **Total** | **55** | **29** | **84** |
