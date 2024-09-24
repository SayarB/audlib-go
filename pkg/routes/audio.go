package routes

import (
	"fmt"
	"path"

	"github.com/gofiber/fiber/v2"

	"github.com/google/uuid"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

type StreamTokenResponse struct {
	Token string `json:"token"`
}

type FilePostRequest struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	MIME     string `json:"mime"`
	Key      string `json:"key"`
}

func AudioRoutes(app *fiber.App) {
	app.Get("/audio/upload", func(c *fiber.Ctx) error {
		fileName := c.Query("filename")
		isAuthenticated := c.Locals("isAuthenticated").(bool)
		if !isAuthenticated {
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}
		key, err := uuid.NewV7()
		if err != nil {
			return c.Status(500).JSON(&ErrorResponse{Message: "Could not generate key"})
		}
		url, err := utils.GetPresignedUploadUrl("sayarsbucket", key.String()+path.Ext(fileName))
		if err != nil {
			return c.Status(500).JSON(&ErrorResponse{Message: "Could not generate presigned URL"})
		}
		return c.Status(200).JSON(&fiber.Map{"url": url, "key": key.String()})
	})
	app.Post("/audiofile", func(c *fiber.Ctx) error {

		isAuthenticated := c.Locals("isAuthenticated").(bool)
		if !isAuthenticated {
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}
		user, err := GetAuthenticatedUser(c)
		if err != nil {
			fmt.Println("Cannot Authenticate User")
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}

		fmt.Println(user.Name)

		fmt.Println("Recieved file")
		// fileHeader, err:=c.FormFile("audioFile")

		fileHeader := &FilePostRequest{}
		err = c.BodyParser(fileHeader)

		if err != nil {
			fmt.Println(err)
			return c.Status(400).JSON(&ErrorResponse{Message: "Could not get file"})
		}

		audioFile := models.AudioFile{
			Key:       fileHeader.Key,
			Folder:    "audio",
			Size:      fileHeader.Size,
			Extension: path.Ext(fileHeader.Filename),
			BucketId:  "sayarsbucket",
			MIMEType:  fileHeader.MIME,
			AuthorId:  user.ID,
		}
		err = database.CreateAudioFile(&audioFile)
		if err != nil {
			return c.Status(500).JSON(&ErrorResponse{Message: "Could not create audio file"})
		}
		return c.Status(201).JSON(&audioFile)
	})
	app.Post("/projectfile", func(c *fiber.Ctx) error {
		isAuthenticated := c.Locals("isAuthenticated").(bool)
		if !isAuthenticated {
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}
		user, err := GetAuthenticatedUser(c)
		if err != nil {
			fmt.Println("Cannot Authenticate User")
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}
		fmt.Println("Recieved file")
		// fileHeader, err:=c.FormFile("audioFile")
		fileHeader := &FilePostRequest{}
		err = c.BodyParser(fileHeader)
		if err != nil {
			fmt.Println(err)
			return c.Status(400).JSON(&ErrorResponse{Message: "Could not get file"})
		}
		projectFile := models.ProjectFile{
			Key:       fileHeader.Key,
			Folder:    "project",
			Size:      fileHeader.Size,
			Extension: path.Ext(fileHeader.Filename),
			BucketId:  "sayarsbucket",
			MIMEType:  fileHeader.MIME,
			AuthorId:  user.ID,
		}
		err = database.CreateProjectFile(&projectFile)
		if err != nil {
			return c.Status(500).JSON(&ErrorResponse{Message: "Could not create project file"})
		}
		return c.Status(201).JSON(&projectFile)
	})
}

// package routes

// import (
// 	"fmt"
// 	"log"
// 	"path"

// 	"github.com/gofiber/fiber/v2"

// 	"github.com/google/uuid"
// 	"github.com/sayar/go-streaming/pkg/models"
// 	"github.com/sayar/go-streaming/pkg/utils"
// 	"github.com/sayar/go-streaming/pkg/utils/database"
// )

// type StreamTokenResponse struct{
// 	Token string `json:"token"`
// }

// func AudioRoutes(app *fiber.App){
// 	app.Post("/audio", func(c *fiber.Ctx) error{

// 		isAuthenticated:=c.Locals("isAuthenticated").(bool)
// 		if !isAuthenticated{
// 			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
// 		}
// 		user,err:=GetAuthenticatedUser(c)
// 		if err!=nil{
// 			fmt.Println("Cannot Authenticate User")
// 			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
// 		}
// 		org, err:=GetCurrentOrganization(c)
// 		if err!=nil{
// 			fmt.Println("organization not found")
// 			return c.Status(401).JSON(&ErrorResponse{Message: "Org not found"})
// 		}

// 		fmt.Print(user.Name)

// 		fmt.Println("Recieved file")
// 		fileHeader, err:=c.FormFile("audioFile")

// 		if err!=nil{
// 			fmt.Println(err)
// 			return c.Status(400).JSON(&ErrorResponse{Message: "Could not get file"})
// 		}

// 		fmt.Println(org.ID)

// 		fmt.Println(fileHeader.Filename)
// 		file,_:=fileHeader.Open()
// 		key,err:=uuid.NewV7()
// 		if err!=nil{
// 			log.Println(err)
// 		}
// 		audioFile:=models.AudioFile{
// 			Key: key.String(),
// 			File: file,
// 			Folder: "audio",
// 			Size: fileHeader.Size,
// 			Extension: path.Ext(fileHeader.Filename),
// 			BucketId: "sayarsbucket",
// 			MIMEType: fileHeader.Header.Get("Content-Type"),
// 			AuthorId: user.ID,
// 		}
// 		_, err=utils.UploadToS3(&audioFile)
// 		if err!=nil{
// 			return c.Status(500).JSON(&ErrorResponse{Message: "Could not upload file to S3"})
// 		}
// 		database.CreateAudioFile(&audioFile)
	
// 		return c.Status(201).JSON(&audioFile)
// 	})
// }