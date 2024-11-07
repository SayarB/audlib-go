package routes

import (
	"context"
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2/organization"
	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/models"
	clerkutils "github.com/sayar/go-streaming/pkg/utils/clerk"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

type OrgChangeRequest struct {
	OrganizationId string `json:"organizationId"`
}

type OrgCreateRequest struct {
	Name string `json:"name"`
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
	app.Get("/orgs/current", func(c *fiber.Ctx) error {
		org, err := GetCurrentOrganization(c)
		if err != nil {
			return c.Status(400).JSON(&ErrorResponse{Message: "cannot get current organization"})
		}
		return c.Status(200).JSON(org)
	})
	app.Post("/orgs", func(c *fiber.Ctx) error {
		user, err := GetAuthenticatedUser(c)
		if err != nil {
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}

		body := &OrgCreateRequest{}
		c.BodyParser(body)

		orgClerkId, err := clerkutils.CreateNewOrganization(user.ClerkId, body.Name)

		if err != nil {
			fmt.Printf("Error creating organization: %v", err)
			return c.Status(500).JSON(&ErrorResponse{Message: "Error creating organization"})
		}
		newOrg := &models.Organization{
			Name:    body.Name,
			ClerkId: orgClerkId,
		}

		err = database.CreateOrganization(newOrg)
		if err != nil {
			fmt.Printf("could not create")
		}

		userOrg := &models.UserOrganization{
			UserId:         user.ID,
			OrganizationId: newOrg.ID,
		}
		err = database.CreateUserOrganization(userOrg)

		if err != nil {
			fmt.Printf("could not create user org")
		}

		return c.Status(201).JSON(&fiber.Map{"message": "Organization created successfully", "clerk_id": newOrg.ClerkId})
	})
	app.Delete("/orgs/:orgId", func(c *fiber.Ctx) error {
		user, err := GetAuthenticatedUser(c)
		if err != nil {
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}
		orgId := c.Params("orgId")
		org := &models.Organization{}
		err = database.GetOrganizationById(&orgId, org)
		if err != nil {
			return c.Status(404).JSON(&ErrorResponse{Message: "Organization not found"})
		}

		userOrg := &models.UserOrganization{
			UserId:         user.ID,
			OrganizationId: org.ID,
		}
		err = database.GetUserOrganization(userOrg)
		if err != nil {
			return c.Status(404).JSON(&ErrorResponse{Message: "User not part of the organization"})
		}
		_, err = organization.Delete(context.Background(), org.ClerkId)
		if err != nil {
			fmt.Println("not able to delete from clerk")
			return c.Status(500).JSON(&ErrorResponse{Message: "Error deleting organization"})
		}
		err = database.DeleteUserOrganization(userOrg.ID)
		if err != nil {
			fmt.Printf("deleted from clerk but user-org not deleted from db : %s", userOrg.ID)
			return c.Status(500).JSON(&ErrorResponse{Message: "Error deleting user organization"})
		}
		err = database.DeleteOrganization(org.ID)
		if err != nil {
			return c.Status(500).JSON(&ErrorResponse{Message: "Error deleting organization"})
		}

		return c.Status(200).JSON(&fiber.Map{"message": "Organization deleted successfully"})
	})
}
