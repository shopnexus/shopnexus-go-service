package model

// PaginationParams represents the pagination parameters
type PaginationParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

// PaginateResult represents a paginated result set
type PaginateResult[T any] struct {
	Data       []T    `json:"data"`
	Total      int64  `json:"total"`
	NextPage   *int32 `json:"next_page,omitempty"`
	NextCursor any    `json:"next_cursor,omitempty"`
}

// NextPage calculates the next page number for pagination
func NextPage(offet, limit int32, total int64) *int32 {
	page := offet + 1
	nextPage := page + 1
	if int64(nextPage*limit) >= total {
		return nil
	}

	return &nextPage
}
