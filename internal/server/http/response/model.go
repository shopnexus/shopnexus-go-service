package response

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CommonResponse struct {
	Data  any    `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
}

type PaginateResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
	Error      *Error     `json:"errors,omitempty"`
}

type Pagination struct {
	Limit      int32   `json:"limit"`
	Page       int32   `json:"page"`
	Total      int64   `json:"total"`
	NextPage   *int32  `json:"next_page"`
	NextCursor *string `json:"next_cursor"`
}
