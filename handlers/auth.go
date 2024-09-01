package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joshuandeleva/go-ticket-backend/models"
	"github.com/go-playground/validator/v10"

)

type AuthHanndler struct {
	service models.AuthService
}

var validate  =  validator.New()


func (r *AuthHanndler) RegisterUser(ctx *fiber.Ctx) error {
	registerData := &models.AuthCredentials{}
	context , cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ctx.BodyParser(registerData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status": "fail",
			"data": nil,
		})
	}
	if err := validate.Struct(registerData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status": "fail",
		})
	}
	token, user, err := r.service.RegisterUser(context, registerData)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status": "fail",
			"data": nil,
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created successfully",
		"status": "success",
		"data": fiber.Map{
			"user": user,
			"token": token,
		},
	})
}

func (r *AuthHanndler) LoginUser(ctx *fiber.Ctx) error {
	loginData := &models.AuthCredentials{}

	context , cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := ctx.BodyParser(&loginData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status": "fail",
			"data": nil,
		})
	}
	// if err := validate.Struct(loginData.Email); err != nil {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"message": err.Error(),
	// 		"status": "fail",
	// 		"data": nil,
	// 	})
	// }

	token, user, err := r.service.LoginUser(context, loginData)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status": "fail",
			"data": nil,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":&fiber.Map{
			"user": user,
			"token": token,
		},
	})
}

func NewAuthHandler(router fiber.Router, service models.AuthService) {
	handler := &AuthHanndler{
		service: service,
	}
	router.Post("/register", handler.RegisterUser)
	router.Post("/login", handler.LoginUser)
}