package model

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type User struct {
	ID                int       `json:"id"`
	Email             string    `json:"email" validate:"required,email"`
	Password          string    `json:"password,omitempty" validate:"required,min=8"`
	EncryptedPassword string    `json:"-"`
	IsAdmin           bool      `json:"is_admin"`
	Banned            bool      `json:"banned"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := hashPassword(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = string(enc)
	}
	return nil
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *User) Sanitize() {
	u.Password = ""
	u.EncryptedPassword = ""
}

func (u *User) ComparePassword(password string) bool {
	return compareHashedPassword(u.EncryptedPassword, password)
}
