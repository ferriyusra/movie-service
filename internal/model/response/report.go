package response

import "github.com/google/uuid"

// RevenueReportItem represents a single row in the revenue report
type RevenueReportItem struct {
	MovieTitle    string  `json:"movieTitle"`
	ShowtimeCount int     `json:"showtimeCount"`
	TicketsSold   int     `json:"ticketsSold"`
	TotalRevenue  float64 `json:"totalRevenue"`
}

// CapacityReportItem represents a single row in the capacity report
type CapacityReportItem struct {
	ShowtimeID     uuid.UUID `json:"showtimeId"`
	MovieTitle     string    `json:"movieTitle"`
	TheaterName    string    `json:"theaterName"`
	TotalSeats     int       `json:"totalSeats"`
	SeatsReserved  int       `json:"seatsReserved"`
	SeatsAvailable int       `json:"seatsAvailable"`
	OccupancyRate  float64   `json:"occupancyRate"`
}

// DashboardSummary represents the admin dashboard summary
type DashboardSummary struct {
	TotalMovies             int            `json:"totalMovies"`
	TotalShowtimes          int            `json:"totalShowtimes"`
	TotalUsers              int            `json:"totalUsers"`
	ActiveReservationsToday int            `json:"activeReservationsToday"`
	RevenueThisWeek         float64        `json:"revenueThisWeek"`
	RevenueLastWeek         float64        `json:"revenueLastWeek"`
	TopMoviesByTickets      []TopMovieItem `json:"topMoviesByTickets"`
	TopMoviesByRevenue      []TopMovieItem `json:"topMoviesByRevenue"`
}

// TopMovieItem represents a movie in the top-N lists
type TopMovieItem struct {
	MovieTitle   string  `json:"movieTitle"`
	TicketsSold  int     `json:"ticketsSold"`
	TotalRevenue float64 `json:"totalRevenue"`
}
