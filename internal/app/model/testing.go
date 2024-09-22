package model

import (
	"database/sql"
	"testing"
)

func TestUser(t *testing.T) *User {
	return &User{
		IDTelegram:  sql.NullInt64{Valid: false},
		Email:       sql.NullString{String: "user@example.org", Valid: true},
		Password:    "password",
		Height:      180,
		Age:         25,
		Weight:      80,
		Gender:      "male",
		PhoneNumber: sql.NullString{Valid: false},
	}
}

func TestUserWithTelegram(t *testing.T) *User {
	return &User{
		IDTelegram:  sql.NullInt64{Int64: 12345678, Valid: true},
		Email:       sql.NullString{Valid: false},
		Password:    "password",
		Height:      160,
		Age:         18,
		Weight:      50,
		Gender:      "female",
		PhoneNumber: sql.NullString{String: "88005553535", Valid: true},
	}
}
