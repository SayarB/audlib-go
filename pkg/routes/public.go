package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

type PublicVersionResponse struct {
	VersionName string
	AuthorName  string
	ProjectName string
}

func PublicRoutes(app *fiber.App) {
	app.Get("/public/version/:versionId/info", func(c *fiber.Ctx) error {
		versionId := c.Params("versionId")
		version, err := database.GetVersionByIdWithProject(versionId)
		if err != nil {
			return c.Status(400).JSON(&fiber.Map{"message": "Invalid version"})
		}
		res := &PublicVersionResponse{
			VersionName: version.Title,
			ProjectName: version.Project.Name,
			AuthorName:  version.Author.DisplayName,
		}
		return c.JSON(res)
	})
}
