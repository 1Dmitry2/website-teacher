package model

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type Admin struct {
	ID                int        `json:"id"`
	Email             string     `json:"email" validate:"required,email"`
	Password          string     `json:"password,omitempty" validate:"omitempty,min=8"`
	PasswordHash      string     `json:"-"`
	ResetToken        *string    `json:"-"`
	ResetTokenExpires *time.Time `json:"-"`
	CreatedAt         time.Time  `json:"created_at"`
}

func (a *Admin) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}

func (a *Admin) BeforeCreate() error {
	if a.Password == "" {
		return errors.New("password is required")
	}
	return a.SetPassword(a.Password)
}

func (a *Admin) SetPassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}
	a.PasswordHash = hash
	a.Password = ""
	return nil
}

func (a *Admin) ComparePassword(password string) bool {
	return compareHashedPassword(a.PasswordHash, password)
}

func (a *Admin) Sanitize() {
	a.Password = ""
	a.PasswordHash = ""
	if a.ResetToken != nil {
		a.ResetToken = nil
	}
	if a.ResetTokenExpires != nil {
		a.ResetTokenExpires = nil
	}
}
