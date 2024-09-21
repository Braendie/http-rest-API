package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID                int    `json:"id"`
	IDTelegram        int    `json:"id_telegram"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Height            int    `json:"height"`
	Age               int    `json:"age"`
	Weight            int    `json:"weight"`
	Gender            string `json:"gender"`
	PhoneNumber       string `json:"phone_number"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.IDTelegram, validation.By(requiredIf(u.Email == ""))),
		validation.Field(&u.Email, validation.By(requiredIf(u.IDTelegram == -1)), is.Email),
		validation.Field(&u.Password,
			validation.By(requiredIf(u.EncryptedPassword == "")),
			validation.By(requiredIf(u.Email != "")),
			validation.Length(6, 30),
		),
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Height, validation.Required),
		validation.Field(&u.Age, validation.Required),
		validation.Field(&u.Weight, validation.Required),
		validation.Field(&u.Gender, validation.Required),
		validation.Field(&u.PhoneNumber),
	)
}

// BeforeCreate ...
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptedString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}
	return nil
}

func (u *User) Sanitize() {
	u.Password = ""
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptedString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
