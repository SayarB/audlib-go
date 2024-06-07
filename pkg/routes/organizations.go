package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

type OrgChangeRequest struct {
	OrganizationId string `json:"organizationId"`
}

func OrganizationRoutes(app *fiber.App) {
	app.Get("/orgs", func(c *fiber.Ctx) error {
		user, err := GetAuthenticatedUser(c)
		if err != nil {
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}

		var orgs []models.UserOrganization

		orgs, err = database.GetUserOrganizationsForUser(user.ID)

		if err != nil {
			return c.Status(500).JSON(&ErrorResponse{Message: "Error fetching organizations"})
		}

		fmt.Println(len(orgs))

		return c.JSON(orgs)
	})
	app.Get("/orgs/check", func(c *fiber.Ctx) error {
		org, err := GetCurrentOrganization(c)
		if err != nil {
			return c.Status(400).JSON(&ErrorResponse{Message: "no organization selected"})
		}
		return c.Status(200).JSON(org)
	})
	app.Get("/orgs/current", func(c *fiber.Ctx) error{
		org, err:= GetCurrentOrganization(c)
		if err!=nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "cannot get current organization"})
		}
		return c.Status(200).JSON(org)
	})
	// app.Post("/orgs/select", func(c *fiber.Ctx) error {

		
	// })
}
