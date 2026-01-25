package common

// struct to be used for any [GET]All APIs that have pagination capabilities
type PaginatedResponse struct {
	Data       any   `json:"data"` // data represents the actual domain model we want to respond with
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
