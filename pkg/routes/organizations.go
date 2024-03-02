package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

type OrgChangeRequest struct{
	OrganizationId string `json:"organizationId"`
}

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
	app.Get("/orgs/check", func(c *fiber.Ctx) error {
		user, err:= GetAuthenticatedUser(c)
		if err!=nil{
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}

		var orgs []models.UserOrganization

		orgs, err = database.GetOrganizationsForUser(user.ID)

		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error fetching organizations"})
		}

		for _, org := range orgs {
			if org.OrganizationId == c.Get("Organization-ID") {
				return c.Status(200).JSON(&fiber.Map{"message": "User is part of the organization"})
			}
		}
		return c.Status(401).JSON(&ErrorResponse{Message: "User is not part of the organization"})
	})
	app.Get("/orgs/current", func(c *fiber.Ctx) error{
		org, err:= GetCurrentOrganization(c)
		if err!=nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "cannot get current organization"})
		}
		return c.Status(200).JSON(org)
	})
	app.Post("/orgs/select", func(c *fiber.Ctx) error {

		token:=c.Cookies("audlib")

		body:=&OrgChangeRequest{}
		err:=c.BodyParser(body)

		if err!=nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "organization ID is not provided"})
		}

		fmt.Println("Token: ", token)
		if token==""{
			return c.Status(400).JSON(&ErrorResponse{Message: "token not provided"})
		}

		fmt.Println("OrganizationId sent = ", body.OrganizationId)

		err=database.ChangeOrganization(token, body.OrganizationId)
		if err!=nil{
			fmt.Println(err)
			return c.Status(500).JSON(&ErrorResponse{Message: "cannot change organization"})
		}
		return c.SendStatus(200)
	})
}