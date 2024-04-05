package routes

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SetSessionCookie(c *fiber.Ctx, token string){
	c.Cookie(&fiber.Cookie{
		Name: "audlib",
		Value: token,
		Expires: time.Now().Add(time.Hour * 24 * 3),
		HTTPOnly: true,
		SameSite: "None",
		Secure: true,
	})
	
}

func AuthRoutes(app *fiber.App) {
	app.Get("/auth/check", func(c *fiber.Ctx) error{
		_, err:= GetAuthenticatedUser(c)
		if err!=nil{
			return c.Status(401).JSON(&fiber.Map{"message":"Unauthorized"})
		}
		return c.SendStatus(200)
	})
	app.Post("/auth/login", func(c *fiber.Ctx) error {
		body := LoginRequest{}
		err:=c.BodyParser(&body)
		if err!=nil{
			fmt.Println(err)
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid request body"})
		}

		// Check if the user exists in the database
		user:=&models.User{}
		err=database.GetUserByEmail(body.Email, user)
		if err!=nil{
			return c.Status(401).JSON(&ErrorResponse{Message: "Invalid credentials"})
		}

		// Check if the password is correct
		if body.Password!=user.Password{
			return c.Status(401).JSON(&ErrorResponse{Message: "Invalid credentials"})
		}

		// Create a session
		token,err:=database.CreateSessionToken(user.ID,nil)
	
		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error creating session"})
		}
		SetSessionCookie(c, token)
		return c.JSON(&user)

	})

	app.Post("/auth/register", func(c *fiber.Ctx) error {
		body := RegisterRequest{}
		err:=c.BodyParser(&body)
		if err!=nil{
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid request body"})
		}
		
		user:=&models.User{
			Name: body.Name,
			Email: body.Email,
			Password: body.Password,
		}
		err=database.CreateUser(user)
		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error creating user"})
		}
		org:=&models.Organization{
			Name: fmt.Sprintf("%s's Org",strings.Split(body.Name, " ")[0]),
		}
		err= database.CreateOrganization(org)
		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error creating user"})
		}
		userOrg:=&models.UserOrganization{
			UserId: user.ID,
			OrganizationId: org.ID,
		}
		err=database.CreateUserOrganization(userOrg)
		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error creating user"})
		}
		token,err:=database.CreateSessionToken(user.ID, nil)

		if err!=nil{
			return c.Status(500).JSON(&ErrorResponse{Message: "Error creating session"})
		}

		SetSessionCookie(c, token)

		return c.JSON(&user)
	})

	app.Get("/auth/info", func(c *fiber.Ctx) error {
		user, err:=GetAuthenticatedUser(c)
		if err!=nil{
			return c.Status(401).JSON(&ErrorResponse{Message: "No Authenticated user"})
		}
		return c.JSON(&user)
	})
	app.Post("/auth/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("session")
		return c.SendStatus(200)
	})
}