package routes

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

type PostProjectRequest struct{
	Name string `json:"name"`
}

func GetAuthenticatedUser(c *fiber.Ctx) (*models.User, error){
	if c.Locals("user")==nil{
		return nil, errors.New("Cannot Get User")
	}
	user:=c.Locals("user").(*models.User)
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
	app.Get("/projects", func(c *fiber.Ctx) error{
		isAuth:=c.Locals("isAuthenticated").(bool)

		if !isAuth{
			fmt.Println("Not Authenticated")
			return c.Status(401).JSON(&ErrorResponse{Message: "Not Authenticated"})
		}
		org, err:=GetCurrentOrganization(c)
		if err!=nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Organization ID"})
		}
		projects,err:=database.GetProjectsByOrganizationId(org.ID)
		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error fetching projects"})
		}
		return c.JSON(&projects)
	})
}