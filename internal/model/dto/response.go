package dto

type Pagination struct {
	Total       int64 `json:"total_records"`
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
}
