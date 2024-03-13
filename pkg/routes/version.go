package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

func VersionRoutes(app *fiber.App) {
	app.Get("/version/:versionId", func(c *fiber.Ctx) error {
		versionId:= c.Params("versionId")

		version, err:=database.GetVersionByIdWithProject(versionId)
		if err!=nil{
			return c.Status(400).JSON(&fiber.Map{"message":"Invalid version"})
		}
		return c.JSON(version)
	})
}