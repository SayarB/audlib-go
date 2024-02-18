package main

import (
	"fmt"
	"log"
	"path"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sayar/go-streaming/pkg/utils"
)

func main() {
	err:=godotenv.Load(".env")
	if err!=nil{
		panic("env cannot be loaded")
	}
	app := fiber.New(fiber.Config{DisablePreParseMultipartForm: true, StreamRequestBody: true})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"hello": "world"})
	})
	app.Post("/audio", func(c *fiber.Ctx) error{
		fmt.Println("Recieved file")
		fmt.Println("Recieved file")
		fileHeader, _:=c.FormFile("audioFile")
		fmt.Println(fileHeader.Filename)
		file,_:=fileHeader.Open()
		key,err:=uuid.NewV7()
		if err!=nil{
			log.Println(err)
		}
		utils.UploadToSupabase(&utils.AudioFile{
			Key: key.String(),
			File: file,
			Folder: "audio",
			Extension: path.Ext(fileHeader.Filename),
			BucketId: "audio",
			MIMEType: fileHeader.Header.Get("Content-Type"),
		})

		return c.JSON(fiber.Map{ "key": key.String()})
	})
	log.Fatal(app.Listen(":3000"))
}
