package routes

import (
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

		orgs:=[]models.UserOrganization{}
		err= database.GetOrganizationsForUser(user.ID, orgs)

		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error fetching organizations"})
		}

		return c.JSON(orgs)
	})
}