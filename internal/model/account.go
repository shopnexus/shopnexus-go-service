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
	ID       []byte `json:"id"` /* unique */
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (AccountBase) ImplementsAccount() {}
func (a AccountBase) GetBase() AccountBase {
	return a
}

type AccountUser struct {
	AccountBase
	Email            string `json:"email"` /* unique */
	Phone            string `json:"phone"` /* unique */
	Gender           int32  `json:"gender"`
	FullName         string `json:"full_name"`
	DefaultAddressID []byte `json:"default_address_id"`
}

type AccountAdmin struct {
	AccountBase
}

type ItemQuantity interface {
	GetID() []byte
	GetQuantity() int64
}

type ItemQuantityBase struct {
	ItemID   []byte `json:"item_id"`
	Quantity int64  `json:"quantity"`
}

func (i ItemQuantityBase) GetID() []byte {
	return i.ItemID
}

func (i ItemQuantityBase) GetQuantity() int64 {
	return i.Quantity
}

type ItemOnCart struct {
	ItemQuantityBase
	CartID []byte `json:"cart_id"`
}

type Cart struct {
	ID       []byte         `json:"id"` /* unique */
	Products []ItemQuantity `json:"products"`
}
