package models

import (
	"context"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type AuthCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstname" validate:"required"`
	LastName string `json:"lastname" validate:"required"`

}



type AuthRepository interface {
	RegisterUser (ctx context.Context, registerData *AuthCredentials) (*User, error)
	GetUser (ctx context.Context, query interface{} , args ...interface{}) (*User, error)
}

type AuthService interface {
	RegisterUser (ctx context.Context, registerData *AuthCredentials) (string , *User, error)
	LoginUser (ctx context.Context, loginData *AuthCredentials) (string ,*User, error)
}

// check if a password matches hash

func MatchesHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// check if email valid

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
