package model

type (
	Role   string
	Gender string
)

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"

	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
	GenderOther  Gender = "OTHER"
)

type Account interface {
	ImplementsAccount()
	GetBase() AccountBase
}

type AccountBase struct {
	ID       int64  `json:"id"` /* unique */
	Username string `json:"username"`
	Role     Role   `json:"role"`
}

func (AccountBase) ImplementsAccount() {}
func (a AccountBase) GetBase() AccountBase {
	return a
}

type AccountUser struct {
	AccountBase
	Email            string `json:"email"` /* unique */
	Phone            string `json:"phone"` /* unique */
	Gender           Gender `json:"gender"`
	FullName         string `json:"full_name"`
	DefaultAddressID int64  `json:"default_address_id"`
}

type AccountAdmin struct {
	AccountBase
}

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

type ItemOnCart struct {
	ItemQuantityBase[int64]
	CartID int64 `json:"cart_id"`
}

type Cart struct {
	ID       int64                 `json:"id"` /* unique */
	Products []ItemQuantity[int64] `json:"products"`
}
