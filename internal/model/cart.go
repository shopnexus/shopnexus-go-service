package model

type ItemOnCart struct {
	ItemQuantityBase[int64]
	CartID int64 `json:"cart_id"`
}

type Cart struct {
	ID       int64                 `json:"id"` /* unique */
	Products []ItemQuantity[int64] `json:"products"`
}
