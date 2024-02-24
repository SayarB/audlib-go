package routes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sayar/go-streaming/pkg/utils"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

func StreamRoutes(app *fiber.App){
	app.Post("/stream/:key/token", func(c *fiber.Ctx) error{
		fmt.Println("Recieved request for File: ", c.Params("key"))
		key:=c.Params("key")
		fileInfo,err:=database.GetFileInfo(key)
		if err!=nil{
			c.Status(500).JSON(&ErrorResponse{Message: "Could not find audio file in db"})
		}

		if isPermitted:=utils.HasPermissionToStream(fileInfo); !isPermitted{
			return c.Status(401).JSON(&ErrorResponse{Message: "No Permissions"})
		}


		if err!=nil{
			c.Status(500).JSON(&ErrorResponse{Message: "Could not find audio file in db"})
		}
		token, err:=utils.GenerateStreamToken(key)
		if err!=nil{
			c.Status(500).JSON(&ErrorResponse{Message: "Could not generate stream token"})
		}
		return c.JSON(&StreamTokenResponse{Token: token})
	})

	app.Get("/stream/:key", func(c *fiber.Ctx) error{
		fmt.Println("Recieved request for File: ", c.Params("key"))
		r:=c.Get("Range")
		streamToken:=c.Query("token")

		token, err:=utils.ValidateStreamToken(streamToken)
		if err!=nil{
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized, Stream Token Invalid"})
		}
		
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claims["key"] != c.Params("key") {
				return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized, Stream Token Invalid"})
			}
		} else {
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized, Stream Token Invalid"})
		}

		fmt.Println(r)

		start := ""
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
		
		key:=c.Params("key")
		fileInfo,err:=database.GetFileInfo(key)
		if err!=nil{
			c.Status(500).JSON(&ErrorResponse{Message: "Could not find audio file in db"})
		}
		result, err:=utils.DownloadFileFromS3(fileInfo, startInt, endInt)
		if err!=nil{
			c.Status(500).JSON(&ErrorResponse{Message: "Could not download file from S3"})
		}

		c.Set("Content-Type", fileInfo.MIMEType)
		c.Set("Content-Length", fmt.Sprintf("%d", len(result)))
		c.Set("Accept-Ranges", "bytes")
		c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", startInt, endInt, fileInfo.Size))


		return c.Status(206).SendString(result)
	})

}