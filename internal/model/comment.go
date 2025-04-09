package model

type Comment struct {
	ID          int64    `json:"id"`
	AccountID   int64    `json:"account_id"`
	DestID      int64    `json:"dest_id"`
	Body        string   `json:"body"`
	Upvote      int64    `json:"upvote"`
	Downvote    int64    `json:"downvote"`
	Score       int32    `json:"score"`
	DateCreated int64    `json:"date_created"`
	DateUpdated int64    `json:"date_updated"`
	Resources   []string `json:"resources"`
}
