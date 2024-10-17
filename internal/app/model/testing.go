package model

import (
	"database/sql"
	"testing"
)

// TestUser returns a test model with email and password for testing.
func TestUser(t *testing.T) *User {
	return &User{
		IDTelegram:  sql.NullInt64{Valid: false},
		Email:       sql.NullString{String: "user@example.org", Valid: true},
		Password:    "password",
	}
}

// TestUserWithTelegram returns a test mode with telegram id for testing.
func TestUserWithTelegram(t *testing.T) *User {
	return &User{
		IDTelegram:  sql.NullInt64{Int64: 12345678, Valid: true},
		Email:       sql.NullString{Valid: false},
		Password:    "password",
	}
}
