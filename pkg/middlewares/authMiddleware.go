package middlewares

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		

		token:=c.Cookies("session") 
		fmt.Println("Token: ", token)
		if token==""{
			fmt.Println("Token not found")
			c.Locals("isAuthenticated", false)
			return c.Next()
		}

		session:=&models.Session{}
		config.DB.Where("token = ?",token).Preload("User").First(session)

		if session.ExpiresAt<time.Now().Unix(){
			c.ClearCookie("session")
			return c.Status(401).JSON(fiber.Map{"message": "Session expired"})
		}

		if session.ID==""{
			fmt.Println("Session not found")
			c.Locals("isAuthenticated", false)
			c.Next()
		}
		c.Locals("user", session.User)
		c.Locals("isAuthenticated", true)

		orgID := c.Get("Organization-ID")
		if orgID == "" {
			fmt.Println("Organization ID not found")
			return c.Next()
		}
		org:=&models.Organization{}
		err:=database.GetOrganizationById(orgID, org)
		if err!=nil{
			fmt.Println("Organization not found")
			return err
		}
		fmt.Println("Organization: ", org.ID)
		// Set the organization ID in c.Locals
		c.Locals("organization", org)

		//check if the user is a member of the organization
		user:=session.User

		userOrg:=&models.UserOrganization{
			UserId: user.ID,
			OrganizationId: org.ID,
		}

		err=database.GetUserOrganization(userOrg)
		if err!=nil{
			return c.Status(401).JSON(fiber.Map{"message": "Organization not allowed for the user"})
		}

		newToken,err:=database.CreateSessionToken(user.ID)

		if err!=nil{
			return c.Status(500).JSON(fiber.Map{"message": "Error creating session"})
		}

		c.Cookie(&fiber.Cookie{
			Name: "session",
			Value: newToken,
			Expires: time.Now().Add(time.Hour * 24 * 3),
		})

		// Continue to the next middleware or route handler
		return c.Next()
	}
}
