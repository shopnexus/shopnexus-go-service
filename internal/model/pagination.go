package model

type PaginationParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
