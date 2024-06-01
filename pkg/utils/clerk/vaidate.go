package utils

import (
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gofiber/fiber/v2"
)

func ValidateClerkToken(c *fiber.Ctx) error {
	claims, ok := clerk.SessionClaimsFromContext(c.Context())
	if !ok {
		return c.Status(401).JSON(&fiber.Map{"message": "Unauthorized"})
	}

	claims.ValidateWithLeeway(time.Now(), 5*time.Second)
	return nil
}
