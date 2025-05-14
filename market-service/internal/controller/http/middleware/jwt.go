package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}
