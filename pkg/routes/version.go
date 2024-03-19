package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

func VersionRoutes(app *fiber.App) {
	app.Get("/version/:versionId", func(c *fiber.Ctx) error {
		versionId:= c.Params("versionId")
		user, err:=GetAuthenticatedUser(c)
		
		if err!=nil{
			return c.Status(401).JSON(&fiber.Map{"message":"Unauthorized, user not found"})
		}
		org, err:=GetCurrentOrganization(c)
		if err!=nil{
			return c.Status(401).JSON(&fiber.Map{"message":"Unauthorized, not in an organization"})
		}

		version, err:=database.GetVersionByIdWithProject(versionId)
		if err!=nil{
			return c.Status(400).JSON(&fiber.Map{"message":"Invalid version"})
		}

		if user.ID!=version.AuthorId{
			project, err:=database.GetProjectById(version.ProjectId)
			if err!=nil{
				return c.Status(400).JSON(&fiber.Map{"message":"Invalid version"})
			}
			if project.OwnerId!=org.ID{
				return c.Status(401).JSON(&fiber.Map{"message":"Unauthorized organization not matching with the project owner"})
			}
		}
		return c.JSON(version)
	})
}