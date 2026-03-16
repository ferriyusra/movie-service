package response

// HealthStatus represents the health status of the application
type HealthStatus struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}
