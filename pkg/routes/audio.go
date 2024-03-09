package routes

import (
	"fmt"
	"log"
	"path"

	"github.com/gofiber/fiber/v2"

	"github.com/google/uuid"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils"
	"github.com/sayar/go-streaming/pkg/utils/database"
)



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

		fmt.Println(org.ID)

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
			AuthorId: user.ID,
		}
		_, err=utils.UploadToS3(&audioFile)
		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Could not upload file to S3"})
		}
		database.CreateAudioFile(&audioFile)
	
		return c.Status(201).JSON(&audioFile)
	})
	
	// app.Get("/audio/stream", func(c *fiber.Ctx) error{	

	// 	r:=c.Get("Range")
	// 	fmt.Println(r)

	// 	start:=""
	// 	end:=""

	// 	startInt:=0
	// 	endInt:=0

	// 	if r != "" {
	// 		rangeParts := strings.Split(r, "=")
	// 		if len(rangeParts) == 2 {
	// 			rangeValues := strings.Split(rangeParts[1], "-")
	// 			if len(rangeValues) == 2 {
	// 				start = rangeValues[0]
	// 				end = rangeValues[1]
	// 			}
	// 		}
	// 	}

	// 	if res, err := strconv.Atoi(start); err!=nil{
	// 		startInt=0
	// 	}else{
	// 		startInt=res
	// 	}

	// 	if res, err := strconv.Atoi(end); err!=nil{
	// 		endInt=startInt+ChunkSize
	// 	}else{
	// 		endInt=res
	// 	}

		
	// 	file,_:=os.Open("./temp.wav")
		
	// 	var buf bytes.Buffer
	// 	io.Copy(&buf, file)
	// 	asString := buf.String()
	// 	actualEnd:=int(math.Min(float64(endInt), float64(len(asString))))
		
	// 	c.Set("Content-Type", "audio/wav")
	// 	c.Set("Content-Length", fmt.Sprintf("%d", len(asString)))
	// 	c.Set("Accept-Ranges", "bytes")
	// 	c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", startInt, actualEnd, len(asString)))
		

	// 	return c.Status(206).SendString(asString[startInt:actualEnd])
	// })
}

