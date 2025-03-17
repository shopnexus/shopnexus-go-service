package model

type (
	Role   string
	Gender string
)

const (
	RoleAdmin Role = "ADMIN"
	RoleStaff Role = "STAFF"
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
	Password string `json:"-"`
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
	DefaultAddressID *int64 `json:"default_address_id"`
}

type AccountAdmin struct {
	AccountBase
}
