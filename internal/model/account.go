package model

type (
	AccountType string
	Gender      string
)

const (
	AccountTypeAdmin AccountType = "ADMIN"
	AccountTypeUser  AccountType = "STAFF"

	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
	GenderOther  Gender = "OTHER"
)

type Account interface {
	ImplementsAccount()
	GetBase() AccountBase
}

type AccountBase struct {
	ID       int64       `json:"id"` /* unique */
	Username string      `json:"username"`
	Password string      `json:"-"`
	Type     AccountType `json:"type"`
}

func (AccountBase) ImplementsAccount() {}
func (a AccountBase) GetBase() AccountBase {
	return a
}

type AccountUser struct {
	AccountBase
	Email            string  `json:"email"` /* unique */
	Phone            string  `json:"phone"` /* unique */
	Gender           Gender  `json:"gender"`
	FullName         string  `json:"full_name"`
	DefaultAddressID *int64  `json:"default_address_id"`
	AvatarURL        *string `json:"avatar_url"`
}

type AccountAdmin struct {
	AccountBase
}
