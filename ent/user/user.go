// Code generated by entc, DO NOT EDIT.

package user

import (
	"fmt"
	"time"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldNickname holds the string denoting the nickname field in the database.
	FieldNickname = "nickname"
	// FieldCarNo holds the string denoting the car_no field in the database.
	FieldCarNo = "car_no"
	// FieldAvatar holds the string denoting the avatar field in the database.
	FieldAvatar = "avatar"
	// FieldSex holds the string denoting the sex field in the database.
	FieldSex = "sex"
	// FieldCity holds the string denoting the city field in the database.
	FieldCity = "city"
	// FieldProvince holds the string denoting the province field in the database.
	FieldProvince = "province"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldDeletedAt holds the string denoting the deleted_at field in the database.
	FieldDeletedAt = "deleted_at"

	// Table holds the table name of the user in the database.
	Table = "users"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldNickname,
	FieldCarNo,
	FieldAvatar,
	FieldSex,
	FieldCity,
	FieldProvince,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultNickname holds the default value on creation for the nickname field.
	DefaultNickname string
	// DefaultAvatar holds the default value on creation for the avatar field.
	DefaultAvatar string
	// DefaultCreatedAt holds the default value on creation for the created_at field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the updated_at field.
	DefaultUpdatedAt func() time.Time
)

// Sex defines the type for the sex enum field.
type Sex string

// Sex values.
const (
	Sex1 Sex = "1"
	Sex2 Sex = "2"
	Sex3 Sex = "3"
)

func (s Sex) String() string {
	return string(s)
}

// SexValidator is a validator for the "sex" field enum values. It is called by the builders before save.
func SexValidator(s Sex) error {
	switch s {
	case Sex1, Sex2, Sex3:
		return nil
	default:
		return fmt.Errorf("user: invalid enum value for sex field: %q", s)
	}
}
