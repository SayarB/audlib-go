package routes

import (
	"fmt"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/middlewares"
	"github.com/sayar/go-streaming/pkg/models"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func SetupRoutes(app *fiber.App) {
	app.Get("/handle-cron", func(c *fiber.Ctx) error {
		var count int64
		config.DB.Model(&models.User{}).Count(&count)
		fmt.Println(count)
		return c.SendStatus(200)
	})
	PublicRoutes(app)
	StreamRoutes(app)
	app.Use([]string{"/audio", "/projects", "/orgs", "/auth/onboard", "/auth/info", "/version", "/projectfile"}, adaptor.HTTPMiddleware(clerkhttp.RequireHeaderAuthorization()))
	app.Use([]string{"/audio", "/projects", "/orgs", "/auth/check", "/auth/info", "/version", "/projectfile"}, middlewares.AuthMiddleware())
	AuthRoutes(app)
	OrganizationRoutes(app)
	AudioRoutes(app)
	ProjectsRoutes(app)
	VersionRoutes(app)
}
