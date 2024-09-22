package model

import (
	"database/sql"

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
	Height            int            `json:"height"`
	Age               int            `json:"age"`
	Weight            int            `json:"weight"`
	Gender            string         `json:"gender"`
	PhoneNumber       sql.NullString `json:"phone_number"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.IDTelegram, validation.By(validationIf(!u.Email.Valid, validation.Required))),
		validation.Field(&u.Email, validation.By(validationIf(!u.IDTelegram.Valid, validation.Required)), validation.By(validationIf(!u.IDTelegram.Valid, is.Email))),
		validation.Field(&u.Password,
			validation.By(validationIf(u.EncryptedPassword.String == "" && u.Email.Valid, validation.Required)),
			validation.Length(6, 30),
		),
		validation.Field(&u.Height, validation.Required),
		validation.Field(&u.Age, validation.Required),
		validation.Field(&u.Weight, validation.Required),
		validation.Field(&u.Gender, validation.Required, validation.By(validationIf(u.Gender == "male" || u.Gender == "female", validation.Required))),
		validation.Field(&u.PhoneNumber, validation.Length(0, 11)),
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
