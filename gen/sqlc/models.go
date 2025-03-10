// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type AccountGender string

const (
	AccountGenderMALE   AccountGender = "MALE"
	AccountGenderFEMALE AccountGender = "FEMALE"
	AccountGenderOTHER  AccountGender = "OTHER"
)

func (e *AccountGender) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AccountGender(s)
	case string:
		*e = AccountGender(s)
	default:
		return fmt.Errorf("unsupported scan type for AccountGender: %T", src)
	}
	return nil
}

type NullAccountGender struct {
	AccountGender AccountGender
	Valid         bool // Valid is true if AccountGender is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAccountGender) Scan(value interface{}) error {
	if value == nil {
		ns.AccountGender, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AccountGender.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAccountGender) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AccountGender), nil
}

type AccountRole string

const (
	AccountRoleADMIN AccountRole = "ADMIN"
	AccountRoleUSER  AccountRole = "USER"
)

func (e *AccountRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AccountRole(s)
	case string:
		*e = AccountRole(s)
	default:
		return fmt.Errorf("unsupported scan type for AccountRole: %T", src)
	}
	return nil
}

type NullAccountRole struct {
	AccountRole AccountRole
	Valid       bool // Valid is true if AccountRole is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAccountRole) Scan(value interface{}) error {
	if value == nil {
		ns.AccountRole, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AccountRole.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAccountRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AccountRole), nil
}

type PaymentPaymentMethod string

const (
	PaymentPaymentMethodCASH  PaymentPaymentMethod = "CASH"
	PaymentPaymentMethodMOMO  PaymentPaymentMethod = "MOMO"
	PaymentPaymentMethodVNPAY PaymentPaymentMethod = "VNPAY"
)

func (e *PaymentPaymentMethod) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentPaymentMethod(s)
	case string:
		*e = PaymentPaymentMethod(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentPaymentMethod: %T", src)
	}
	return nil
}

type NullPaymentPaymentMethod struct {
	PaymentPaymentMethod PaymentPaymentMethod
	Valid                bool // Valid is true if PaymentPaymentMethod is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentPaymentMethod) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentPaymentMethod, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentPaymentMethod.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentPaymentMethod) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentPaymentMethod), nil
}

type PaymentRefundMethod string

const (
	PaymentRefundMethodDROPOFF PaymentRefundMethod = "DROP_OFF"
	PaymentRefundMethodPICKUP  PaymentRefundMethod = "PICK_UP"
)

func (e *PaymentRefundMethod) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentRefundMethod(s)
	case string:
		*e = PaymentRefundMethod(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentRefundMethod: %T", src)
	}
	return nil
}

type NullPaymentRefundMethod struct {
	PaymentRefundMethod PaymentRefundMethod
	Valid               bool // Valid is true if PaymentRefundMethod is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentRefundMethod) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentRefundMethod, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentRefundMethod.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentRefundMethod) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentRefundMethod), nil
}

type PaymentStatus string

const (
	PaymentStatusPENDING   PaymentStatus = "PENDING"
	PaymentStatusSUCCESS   PaymentStatus = "SUCCESS"
	PaymentStatusCANCELLED PaymentStatus = "CANCELLED"
	PaymentStatusFAILED    PaymentStatus = "FAILED"
)

func (e *PaymentStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentStatus(s)
	case string:
		*e = PaymentStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentStatus: %T", src)
	}
	return nil
}

type NullPaymentStatus struct {
	PaymentStatus PaymentStatus
	Valid         bool // Valid is true if PaymentStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentStatus), nil
}

type AccountAddress struct {
	ID         int64
	UserID     int64
	Address    string
	City       string
	Province   string
	Country    string
	PostalCode string
}

type AccountAdmin struct {
	ID int64
}

type AccountBase struct {
	ID       int64
	Username string
	Password string
	Role     AccountRole
}

type AccountCart struct {
	ID int64
}

type AccountItemOnCart struct {
	CartID         int64
	ProductModelID int64
	Quantity       int64
}

type AccountUser struct {
	ID               int64
	Email            string
	Phone            string
	Gender           AccountGender
	FullName         string
	DefaultAddressID pgtype.Int8
}

type PaymentBase struct {
	ID          int64
	UserID      int64
	Method      PaymentPaymentMethod
	Status      PaymentStatus
	Address     string
	Total       int64
	DateCreated pgtype.Timestamptz
}

type PaymentProductOnPayment struct {
	PaymentID       int64
	ProductSerialID string
	Quantity        int64
	Price           int64
	TotalPrice      int64
}

type PaymentRefund struct {
	ID          int64
	PaymentID   int64
	Method      PaymentRefundMethod
	Status      PaymentStatus
	Reason      string
	Address     pgtype.Text
	DateCreated pgtype.Timestamptz
	DateUpdated pgtype.Timestamptz
}

type ProductBase struct {
	ID             int64
	SerialID       string
	ProductModelID int64
	Sold           bool
	DateCreated    pgtype.Timestamptz
	DateUpdated    pgtype.Timestamptz
}

type ProductBrand struct {
	ID          int64
	Name        string
	Description string
}

type ProductModel struct {
	ID               int64
	BrandID          int64
	Name             string
	Description      string
	ListPrice        int64
	DateManufactured pgtype.Timestamptz
}

type ProductResource struct {
	OwnerID int64
	S3ID    string
}

type ProductSale struct {
	ID              int64
	TagName         pgtype.Text
	ProductModelID  pgtype.Int8
	DateStarted     pgtype.Timestamptz
	DateEnded       pgtype.Timestamptz
	Quantity        int64
	Used            int64
	IsActive        bool
	DiscountPercent pgtype.Int8
	DiscountPrice   pgtype.Int8
}

type ProductTag struct {
	TagName     string
	Description string
}

type ProductTagOnProduct struct {
	ProductModelID int64
	TagName        string
}
