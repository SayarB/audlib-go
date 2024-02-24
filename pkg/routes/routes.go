package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/middlewares"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func SetupRoutes(app *fiber.App){
	StreamRoutes(app)
	app.Use([]string{"/audio","/project"},middlewares.AuthMiddleware())
	AuthRoutes(app)
	AudioRoutes(app)
	ProjectsRoutes(app)
}