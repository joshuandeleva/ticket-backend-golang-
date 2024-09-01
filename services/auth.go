package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joshuandeleva/go-ticket-backend/models"
	"github.com/joshuandeleva/go-ticket-backend/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	repository models.AuthRepository
}

func (s *AuthService) RegisterUser(context context.Context,registerData *models.AuthCredentials) (string ,*models.User, error) {
	if !models.IsValidEmail(registerData.Email) {
		return "", nil, fmt.Errorf("please provide a valid email")
	}

	if _ , err := s.repository.GetUser(context, "email =?", registerData.Email); !errors.Is(err, gorm.ErrRecordNotFound){
		return "", nil, fmt.Errorf("user with email %s already exists", registerData.Email)
	} 

	hashedPassword, err := utils.HashPassword(registerData.Password)
	if err != nil {
		return "", nil, err
	}
	registerData.Password = hashedPassword
	user, err := s.repository.RegisterUser(context, registerData)
	if err != nil {
		return "", nil, err
	}
	claims := jwt.MapClaims{
		"id": user.ID,
		"role": user.Role,
		"exp": time.Now().Add(time.Hour * 100).Unix(),
	}
	token , error := utils.GenerateJWT(claims , jwt.SigningMethodHS256 , os.Getenv("JWT_SECRET"))
	if error != nil {
		return "", nil, error
	}
	return token, user, nil
}

func (s *AuthService) LoginUser(context context.Context , loginData *models.AuthCredentials) (string ,*models.User, error) {
	user, err:= s.repository.GetUser(context, "email =?", loginData.Email)
	fmt.Println(loginData.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, fmt.Errorf("invalid credentials")
		}
		return "", nil, err
	}
	if !models.MatchesHash(loginData.Password, user.Password) {
		return "", nil, fmt.Errorf("password is incorrect")
	}

	claims := jwt.MapClaims{
		"id": user.ID,
		"role": user.Role,
		"exp": time.Now().Add(time.Hour * 100).Unix(),
	}
	token , error := utils.GenerateJWT(claims , jwt.SigningMethodHS256 , os.Getenv("JWT_SECRET"))
	if error != nil {
		return "", nil, fmt.Errorf("unable to generate token")
	}
	return token, user, nil
}

func NewAuthService(repository models.AuthRepository) models.AuthService {
	return &AuthService{
		repository: repository,
	}
}