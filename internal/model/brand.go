package model

type Brand struct {
	ID          int64  `json:"id"` /* unique */
	Name        string `json:"name"`
	Description string `json:"description"`

	Resources []string `json:"resources"`
}
