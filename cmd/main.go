package main

import (
	"log"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/routes"
)

func init() {
	err:=godotenv.Load(".env")
	if err!=nil{
		panic("env cannot be loaded")
	}
	config.ConnectToDatabase()
	config.ConfigureS3Client()
}

func main() {
	app := fiber.New(fiber.Config{DisablePreParseMultipartForm: true, StreamRequestBody: true})

	app.Use(cors.New(cors.Config{
		AllowOrigins:    "*",
	}))
	routes.SetupRoutes(app)
	
	log.Fatal(app.Listen(":8000"))
}
