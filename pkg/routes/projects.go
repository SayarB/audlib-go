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

type CreateVersion struct{
	Title string `json:"title"`
	AudioFileId string `json:"audioFileId"`
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
		if err!=nil || org==nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Organization ID"})
		}
		projects,err:=database.GetProjectsByOrganizationId(org.ID)
		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error fetching projects"})
		}
		return c.JSON(&projects)
	})
	app.Post("/projects/:projectId/version", func(c *fiber.Ctx) error{

		isAuthenticated:=c.Locals("isAuthenticated").(bool)
		if !isAuthenticated{
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}
		user, err:=GetAuthenticatedUser(c)
		if err!=nil{
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}
		org, err:=GetCurrentOrganization(c)

		if err!=nil || org==nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Organization ID"})
		}

		projectId:=c.Params("projectId")

		if projectId==""{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Project ID"})
		}

		project:= &models.Project{}

		project,err=database.GetProjectById(projectId)

		if err!=nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Project ID"})
		}

		if project.OwnerId!=org.ID{
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized. Not Owned by the Organization"})
		}

		

		body:=&CreateVersion{}
		err=c.BodyParser(body)

		if err!=nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Request Body"})
		}

		version:=&models.Version{
			Title: body.Title,
			AudioFileId: body.AudioFileId,
			IsPublished: false,
			AuthorId: user.ID,
			ProjectId: projectId,
		}

		err=database.CreateVersion(version)

		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error creating version"})
		}

		return c.Status(201).JSON(&version)
	})
	app.Get("/projects/:projectId/versions", func(c *fiber.Ctx) error{
		isAuthenticated:=c.Locals("isAuthenticated").(bool)
		if !isAuthenticated{
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}
		org, err:=GetCurrentOrganization(c)

		if err!=nil || org==nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Organization ID"})
		}

		projectId:=c.Params("projectId")

		if projectId==""{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Project ID"})
		}

		project:= &models.Project{}

		project,err=database.GetProjectById(projectId)

		if err!=nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Project ID"})
		}

		if project.OwnerId!=org.ID{
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized. Not Owned by the Organization"})
		}
		versions,err:=database.GetVersionsByProjectId(projectId)
		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error fetching versions"})
		}
		return c.JSON(&versions)
	})
	app.Get("/projects/:projectId/versions/:versionId", func(c *fiber.Ctx) error{
		isAuthenticated:=c.Locals("isAuthenticated").(bool)
		if !isAuthenticated{
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized"})
		}
		org, err:=GetCurrentOrganization(c)

		if err!=nil || org==nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Organization ID"})
		}

		projectId:=c.Params("projectId")
		versionId:=c.Params("versionId")

		if projectId=="" || versionId==""{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Project ID"})
		}

		project:= &models.Project{}

		project,err=database.GetProjectById(projectId)

		if err!=nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid Project ID"})
		}

		if project.OwnerId!=org.ID{
			return c.Status(401).JSON(&ErrorResponse{Message: "Unauthorized. Not Owned by the Organization"})
		}
		version,err:=database.GetVersionById(versionId)
		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error fetching version"})
		}
		return c.JSON(&version)
	})
}