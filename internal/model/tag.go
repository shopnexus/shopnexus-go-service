package model

type Tag struct {
	Name        string `json:"name"` /* unique */
	Description string `json:"description"`
}
