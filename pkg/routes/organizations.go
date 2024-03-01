package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

func OrganizationRoutes(app *fiber.App) {
	app.Get("/orgs", func(c *fiber.Ctx) error {
		user, err:= GetAuthenticatedUser(c)
		if err!=nil{
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}

		var orgs []models.UserOrganization

		orgs, err = database.GetOrganizationsForUser(user.ID)

		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error fetching organizations"})
		}

		fmt.Println(len(orgs))

		return c.JSON(orgs)
	})
}