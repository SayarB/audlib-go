package routes

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/google/uuid"
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

const ChunkSize = 250*1000

type StreamTokenResponse struct{
	Token string `json:"token"`
}

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
			Size: fileHeader.Size,
			Extension: path.Ext(fileHeader.Filename),
			BucketId: "sayarsbucket",
			MIMEType: fileHeader.Header.Get("Content-Type"),
			ProjectId: projectId,
		}
		_, err=utils.UploadToS3(&audioFile)
		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Could not upload file to S3"})
		}
		database.CreateAudioFile(config.DB, &audioFile)
	
		return c.Status(201).JSON(&audioFile)
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

	
	app.Get("/audio/stream", func(c *fiber.Ctx) error{	

		r:=c.Get("Range")
		fmt.Println(r)

		start:=""
		end:=""

		startInt:=0
		endInt:=0

		if r != "" {
			rangeParts := strings.Split(r, "=")
			if len(rangeParts) == 2 {
				rangeValues := strings.Split(rangeParts[1], "-")
				if len(rangeValues) == 2 {
					start = rangeValues[0]
					end = rangeValues[1]
				}
			}
		}

		if res, err := strconv.Atoi(start); err!=nil{
			startInt=0
		}else{
			startInt=res
		}

		if res, err := strconv.Atoi(end); err!=nil{
			endInt=startInt+ChunkSize
		}else{
			endInt=res
		}

		
		file,_:=os.Open("./temp.wav")
		
		var buf bytes.Buffer
		io.Copy(&buf, file)
		asString := buf.String()
		actualEnd:=int(math.Min(float64(endInt), float64(len(asString))))
		
		c.Set("Content-Type", "audio/wav")
		c.Set("Content-Length", fmt.Sprintf("%d", len(asString)))
		c.Set("Accept-Ranges", "bytes")
		c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", startInt, actualEnd, len(asString)))
		

		return c.Status(206).SendString(asString[startInt:actualEnd])
	})
}

