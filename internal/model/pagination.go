package model

type PaginationParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}
