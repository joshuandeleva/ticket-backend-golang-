package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWT(claims jwt.MapClaims , method jwt.SigningMethod , jswtsecret string) (string, error) {
	return jwt.NewWithClaims(method , claims).SignedString([]byte(jswtsecret))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}