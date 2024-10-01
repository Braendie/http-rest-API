package model

import (
	"database/sql"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int            `json:"id"`
	IDTelegram        sql.NullInt64  `json:"id_telegram"`
	Email             sql.NullString `json:"email"`
	Password          string         `json:"password,omitempty"`
	EncryptedPassword sql.NullString `json:"-"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.IDTelegram, validation.By(validationIf(!u.Email.Valid, validation.Required)), validation.By(validationIf(!u.Email.Valid, validation.Min(0)))),
		validation.Field(&u.Email, validation.By(validationIf(!u.IDTelegram.Valid, validation.Required)), validation.By(validationIf(!u.IDTelegram.Valid, is.Email))),
		validation.Field(&u.Password,
			validation.By(validationIf(u.EncryptedPassword.String == "" && u.Email.Valid, validation.Required)),
			validation.Length(6, 30),
		),
	)
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptedString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword.String = enc
		u.EncryptedPassword.Valid = true
	}
	return nil
}

func (u *User) Sanitize() {
	u.Password = ""
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword.String), []byte(password)) == nil
}

func encryptedString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func CheckPassword(password string) bool {
	if len(password) < 6 {
		return false
	} else if password == strings.ToUpper(strings.ToLower(password)) || password == strings.ToLower(strings.ToUpper(password)) {
		return false
	} else if !strings.ContainsAny(password, "1234567890") {
		return false
	}
	return true
}
