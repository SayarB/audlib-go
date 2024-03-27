package routes

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sayar/go-streaming/pkg/utils"
	"github.com/sayar/go-streaming/pkg/utils/database"
)
const ChunkSize = 250*1000
func StreamRoutes(app *fiber.App){
	app.Post("/stream/:id/token", func(c *fiber.Ctx) error{
		fmt.Println("Stream Token Request for Version : ", c.Params("id"))
		Vid:=c.Params("id")
		version,err:=database.GetVersionById(Vid)
		if err!=nil{
			c.Status(500).JSON(&ErrorResponse{Message: "Could not find audio file in db"})
		}

		if isPermitted:=utils.HasPermissionToStream(version); !isPermitted{
			return c.Status(401).JSON(&ErrorResponse{Message: "No Permissions"})
		}


		if err!=nil{
			c.Status(500).JSON(&ErrorResponse{Message: "Could not find audio file in db"})
		}
		token, err:=utils.GenerateStreamToken(version)
		if err!=nil{
			c.Status(500).JSON(&ErrorResponse{Message: "Could not generate stream token"})
		}
		return c.JSON(&StreamTokenResponse{Token: token})
	})

	app.Get("/stream/:id", func(c *fiber.Ctx) error{
		fmt.Println("Recieved request for File: ", c.Params("id"))
		r:=c.Get("Range")
		streamToken:=c.Query("token")

		token, err:=utils.ValidateStreamToken(streamToken)
		if err!=nil{
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized, Stream Token Invalid"})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			if claims["vid"] != c.Params("id") {
				return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized, Stream Token Invalid"})
			}
		} else {
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized, Stream Token Invalid"})
		}

		key:= claims["key"].(string)
		mime:= claims["mime"].(string)
		bucket:= claims["bucket"].(string)
		size:= int(claims["size"].(float64))
		extension:= claims["extension"].(string)
		

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
		
		actualEnd:=int(math.Min(float64(endInt), float64(size)-1))

		result, err:=utils.DownloadFileFromS3(&utils.S3DownloadInput{Key: key, BucketId: bucket, Extension: extension}, startInt, actualEnd)
		if err!=nil{
			c.Status(500).JSON(&ErrorResponse{Message: "Could not download file from S3"})
		}

		c.Set("Content-Type", mime)
		c.Set("Content-Length", fmt.Sprintf("%d", len(result)))
		c.Set("Accept-Ranges", "bytes")
		c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", startInt, actualEnd, size))


		return c.Status(206).SendString(result)
	})

}