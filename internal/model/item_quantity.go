package model

type ItemQuantity[T any] interface {
	GetID() T
	GetQuantity() int64
}

type ItemQuantityBase[T any] struct {
	ItemID   T     `json:"item_id"`
	Quantity int64 `json:"quantity"`
}

func (i ItemQuantityBase[T]) GetID() T {
	return i.ItemID
}

func (i ItemQuantityBase[T]) GetQuantity() int64 {
	return i.Quantity
}
