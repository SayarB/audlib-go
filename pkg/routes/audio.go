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

		isAuthenticated:=c.Locals("isAuthenticated").(bool)
		if !isAuthenticated{
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}
		user,err:=GetAuthenticatedUser(c)
		if err!=nil{
			fmt.Println("Cannot Authenticate User")
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}
		org, err:=GetCurrentOrganization(c)
		if err!=nil{
			fmt.Println("organization not found")
			return c.Status(401).JSON(&ErrorResponse{Message: "Org not found"})
		}

		fmt.Print(user.Name)

		fmt.Println("Recieved file")
		fileHeader, _:=c.FormFile("audioFile")
		projectId :=c.FormValue("projectId")

		project, err:=database.GetProjectById(projectId)
		if err!=nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Project ID"})
		}

		fmt.Println(project.OwnerId)
		fmt.Println(org.ID)

		if project.OwnerId!=org.ID{
			return c.Status(401).JSON(&ErrorResponse{Message: "Project not in organization"})
		}

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

		org:=c.Locals("organization").(*models.Organization)

		if org.ID!=fileInfo.Project.OwnerId{
			return c.Status(401).JSON(&ErrorResponse{Message: "Project is not owned by the organization"})
		}

		

		if err!=nil{
			c.Status(500).JSON(&ErrorResponse{Message: "Could not find audio file in db"})
		}
		return c.JSON(fileInfo)
	})
}
