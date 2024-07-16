package main

import (
	"fmt"
	"log"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/routes"
)

func init() {

	_, err := os.Open(".env")

	if !(err != nil && os.IsNotExist(err)) {
		err := godotenv.Load(".env")
		if err != nil {
			panic("env cannot be loaded")
		}
	}
	_, err = config.ConnectToDatabase()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to database: %v", err))

	}
	config.ConfigureS3Client()
}

func main() {
	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))
	app := fiber.New(fiber.Config{DisablePreParseMultipartForm: true, StreamRequestBody: true})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_ORIGIN"),
		AllowCredentials: true,
	}))

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
