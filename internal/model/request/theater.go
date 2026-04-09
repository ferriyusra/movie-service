package request

// CreateTheaterRequest is the input for creating a theater
type CreateTheaterRequest struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seatsPerRow"`
}

// UpdateTheaterRequest is the input for updating a theater
type UpdateTheaterRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}
