package model

import (
	"errors"
	"net"
	"strings"
	"github.com/go-playground/validator/v10"
	"time"
)

type User struct {
	ID                              int        `json:"id"`
	Email                           string     `json:"email" validate:"required,email"`
	Password                        string     `json:"password,omitempty" validate:"required,min=8"`
	EncryptedPassword               string     `json:"-"`
	IsAdmin                         bool       `json:"is_admin"`
	Banned                          bool       `json:"banned"`
	EmailVerified                   bool       `json:"email_verified"`
	EmailVerificationToken          *string    `json:"-"`
	EmailVerificationTokenExpires   *time.Time `json:"-"`
	ResetToken                      *string    `json:"-"`
	ResetTokenExpires               *time.Time `json:"-"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at"`
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
	if err := validate.Struct(u); err != nil {
		return err
	}
	
	// Дополнительная проверка существования email домена
	if u.Email != "" {
		if err := ValidateEmailDomain(u.Email); err != nil {
			return err
		}
	}
	
	return nil
}

func (u *User) Sanitize() {
	u.Password = ""
	u.EncryptedPassword = ""
}

func (u *User) ComparePassword(password string) bool {
	return compareHashedPassword(u.EncryptedPassword, password)
}

func (u *User) SetPassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	enc, err := hashPassword(password)
	if err != nil {
		return err
	}
	u.EncryptedPassword = string(enc)
	u.Password = ""
	return nil
}

// ValidateEmailDomain проверяет существование MX записей для домена email
func ValidateEmailDomain(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("invalid email format")
	}
	
	domain := parts[1]
	
	// Проверяем MX записи
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		// Если MX записей нет, проверяем A запись (некоторые домены используют A записи для почты)
		_, err := net.LookupHost(domain)
		if err != nil {
			return errors.New("email domain does not exist or is not configured for email")
		}
	}
	
	return nil
}
