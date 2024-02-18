package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

type PostProjectRequest struct{
	Name string `json:"name"`
}

func GetAuthenticatedUser(c *fiber.Ctx) (*models.User, error){
	user:=&models.User{}
	userId:=c.Locals("user").(string)
	err:=database.GetUserById(userId, user)
	if err!=nil{
		return nil,err
	}
	return user,nil
}

func GetCurrentOrganization(c *fiber.Ctx) (*models.Organization, error){
	
	org,ok:=c.Locals("organization").(*models.Organization)
	if !ok{
		return nil, nil
	}
	return org,nil
}

func ProjectsRoutes(app *fiber.App){
	app.Post("/projects", func(c *fiber.Ctx) error{
		body := PostProjectRequest{}
		err:=c.BodyParser(&body)
		if err!=nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid request body"})
		}
		// user,err:=GetAuthenticatedUser(c)
		org, err:=GetCurrentOrganization(c)
		if err!=nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Organization ID"})
		}
		
		projectObj:=&models.Project{
			Name: body.Name,
			OwnerId: org.ID,
		}
		err=database.CreateProject(projectObj)
		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error creating project"})
		}
		return c.Status(201).JSON(&projectObj)
	})
}