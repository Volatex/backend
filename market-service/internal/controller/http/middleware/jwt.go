package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func ExtractUserID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return fiber.ErrUnauthorized
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		c.Locals("user_id", userID)
		return c.Next()
	}
}

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret"), nil // заменишь на свой ключ
		})
		if err != nil || !token.Valid {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["user_id"].(string)

		c.Locals("user_id", userID)

		return c.Next()
	}
}
