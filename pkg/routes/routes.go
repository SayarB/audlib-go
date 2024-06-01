package routes

import (
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/sayar/go-streaming/pkg/middlewares"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func SetupRoutes(app *fiber.App) {
	StreamRoutes(app)
	app.Use([]string{"/audio", "/project", "/orgs", "/auth/onboard", "/auth/info", "/version"}, adaptor.HTTPMiddleware(clerkhttp.RequireHeaderAuthorization()))
	app.Use([]string{"/audio", "/project", "/orgs", "/auth/check", "/auth/info", "/version"}, middlewares.AuthMiddleware())
	AuthRoutes(app)
	OrganizationRoutes(app)
	AudioRoutes(app)
	ProjectsRoutes(app)
	VersionRoutes(app)
}
