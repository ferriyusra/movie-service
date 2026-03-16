package response

// APIResponse is the standardized envelope for all API responses
type APIResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    interface{}       `json:"data,omitempty"`
	Meta    *Meta             `json:"meta,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// Meta holds pagination metadata
type Meta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// OK returns a success response with data
func OK(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// OKWithMeta returns a success response with data and pagination metadata
func OKWithMeta(message string, data interface{}, meta *Meta) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

// Err returns an error response
func Err(message string) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
	}
}

// ValidationErr returns a validation error response with field-level errors
func ValidationErr(message string, errors map[string]string) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}
