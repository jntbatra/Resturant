package response

// PaginatedResponse wraps list responses with pagination metadata
type PaginatedResponse struct {
	Data       interface{}        `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// PaginationMetadata contains pagination information
type PaginationMetadata struct {
	Offset  int  `json:"offset"`
	Limit   int  `json:"limit"`
	Total   int  `json:"total"`
	HasMore bool `json:"has_more"`
}

// NewPaginatedResponse creates a paginated response with metadata
func NewPaginatedResponse(data interface{}, offset int, limit int, total int) *PaginatedResponse {
	hasMore := (offset + limit) < total
	return &PaginatedResponse{
		Data: data,
		Pagination: PaginationMetadata{
			Offset:  offset,
			Limit:   limit,
			Total:   total,
			HasMore: hasMore,
		},
	}
}

// SuccessResponse is used for non-list responses
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}, message string) *SuccessResponse {
	return &SuccessResponse{
		Data:    data,
		Message: message,
	}
}
