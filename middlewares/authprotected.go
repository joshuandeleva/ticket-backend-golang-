package middlewares

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joshuandeleva/go-ticket-backend/models"
	"gorm.io/gorm"
)

func AuthProtected(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			log.Warn("Authorization header is missing")
			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"message": "Unauthorized",
				"status":  "fail",
			})
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Warn("Invalid authorization header")
			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"message": "Unauthorized",
				"status":  "fail",
			})
		}
		tokenstring := tokenParts[1]
		secret := []byte(os.Getenv("JWT_SECRET"))
		token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			log.Warn("Invalid token")
			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"message": "Unauthorized",
				"status":  "fail",
			})
		}
		userdId := token.Claims.(jwt.MapClaims)["id"]
		if err := db.Model(&models.User{}).Where("id = ?", userdId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found")
			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"message": "Unauthorized",
				"status":  "fail",
			})
		}
		ctx.Locals("userId", userdId)
		return ctx.Next()
	}
}
