package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the organization ID from the header
		orgID := c.Get("Organization-ID")
		org:=&models.Organization{}
		err:=database.GetOrganizationById(orgID, org)
		if err!=nil{
			return err
		}
		// Set the organization ID in c.Locals
		c.Locals("organization", org)

		// Continue to the next middleware or route handler
		return c.Next()
	}
}
