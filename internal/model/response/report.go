package response

import "github.com/google/uuid"

// RevenueReportItem represents a single row in the revenue report
type RevenueReportItem struct {
	MovieTitle    string  `json:"movie_title"`
	ShowtimeCount int     `json:"showtime_count"`
	TicketsSold   int     `json:"tickets_sold"`
	TotalRevenue  float64 `json:"total_revenue"`
}

// CapacityReportItem represents a single row in the capacity report
type CapacityReportItem struct {
	ShowtimeID    uuid.UUID `json:"showtime_id"`
	MovieTitle    string  `json:"movie_title"`
	TheaterName   string  `json:"theater_name"`
	TotalSeats    int     `json:"total_seats"`
	SeatsReserved int     `json:"seats_reserved"`
	SeatsAvailable int    `json:"seats_available"`
	OccupancyRate  float64 `json:"occupancy_rate"`
}

// DashboardSummary represents the admin dashboard summary
type DashboardSummary struct {
	TotalMovies             int                `json:"total_movies"`
	TotalShowtimes          int                `json:"total_showtimes"`
	TotalUsers              int                `json:"total_users"`
	ActiveReservationsToday int                `json:"active_reservations_today"`
	RevenueThisWeek         float64            `json:"revenue_this_week"`
	RevenueLastWeek         float64            `json:"revenue_last_week"`
	TopMoviesByTickets      []TopMovieItem     `json:"top_movies_by_tickets"`
	TopMoviesByRevenue      []TopMovieItem     `json:"top_movies_by_revenue"`
}

// TopMovieItem represents a movie in the top-N lists
type TopMovieItem struct {
	MovieTitle   string  `json:"movie_title"`
	TicketsSold  int     `json:"tickets_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}
