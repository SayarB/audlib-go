package routes

import (
	"fmt"
	"log"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

func AudioRoutes(app *fiber.App){
	app.Post("/audio", func(c *fiber.Ctx) error{
		fmt.Println("Recieved file")
		fmt.Println("Recieved file")
		fileHeader, _:=c.FormFile("audioFile")
		projectId :=c.FormValue("projectId")
		fmt.Println(fileHeader.Filename)
		file,_:=fileHeader.Open()
		key,err:=uuid.NewV7()
		if err!=nil{
			log.Println(err)
		}
		audioFile:=models.AudioFile{
			Key: key.String(),
			File: file,
			Folder: "audio",
			Extension: path.Ext(fileHeader.Filename),
			BucketId: "audio",
			MIMEType: fileHeader.Header.Get("Content-Type"),
			ProjectId: projectId,
		}
		utils.UploadToSupabase(&audioFile)
		database.CreateAudioFile(config.DB, &audioFile)
	
		return c.JSON(&audioFile)
	})
	app.Get("/audio/:key/info", func(c *fiber.Ctx) error{
		fmt.Println("Recieved request for File: ", c.Params("key"))
		key:=c.Params("key")
		fileInfo,err:=database.GetFileInfo(key)
		if err!=nil{
			c.Status(500).JSON(&ErrorResponse{Message: "Could not find audio file in db"})
		}
		return c.JSON(fileInfo)
	})
}
