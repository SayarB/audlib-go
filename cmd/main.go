package main

import (
	"log"

	fiber "github.com/gofiber/fiber/v2"
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
}

func main() {
	app := fiber.New(fiber.Config{DisablePreParseMultipartForm: true, StreamRequestBody: true})

	routes.SetupRoutes(app)
	
	log.Fatal(app.Listen(":3000"))
}
