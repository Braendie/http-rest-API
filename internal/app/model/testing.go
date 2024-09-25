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
	}
}

func TestUserWithTelegram(t *testing.T) *User {
	return &User{
		IDTelegram:  sql.NullInt64{Int64: 12345678, Valid: true},
		Email:       sql.NullString{Valid: false},
		Password:    "password",
	}
}
