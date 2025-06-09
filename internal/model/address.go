package model

type Address struct {
	ID          int64  `json:"id"` /* unique */
	UserID      int64  `json:"user_id"`
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Province    string `json:"province"`
	Country     string `json:"country"`
	DateCreated int64  `json:"date_created"`
}
