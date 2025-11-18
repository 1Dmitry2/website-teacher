package apiserver

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type TokenScope string

const (
	TokenScopeUser  TokenScope = "user"
	TokenScopeAdmin TokenScope = "admin"
)

type Claims struct {
	SubjectID int        `json:"subject_id"`
	Scope     TokenScope `json:"scope"`
	jwt.RegisteredClaims
}

func GenerateToken(subjectID int, scope TokenScope, secret string, expiration time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		SubjectID: subjectID,
		Scope:     scope,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(subjectID),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
