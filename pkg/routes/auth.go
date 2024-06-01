package routes

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/organization"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gofiber/fiber/v2"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils/database"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	OrgName     string `json:"orgName"`
}

func SetSessionCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "audlib",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 3),
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
	})

}

func GetAuthorizationToken(c *fiber.Ctx) (string, error) {
	header := c.Get("Authorization")
	if header == "" {
		return "", errors.New("Authorization header not found")
	}
	token := strings.TrimPrefix(header, "Bearer ")
	return token, nil
}

func createNewUser(c *fiber.Ctx, user *models.User) error {
	err := database.CreateUser(user)
	if err != nil {
		return c.Status(500).JSON(&ErrorResponse{Message: "Error creating user"})
	}
	return c.Status(201).JSON(&user)

}

// func sendMagicLink(c *fiber.Ctx, email string) error {
// 	magicLinkToken, err := utils.GenerateMagicLinkToken(email)
// 	if err != nil {
// 		return c.Status(500).JSON(&ErrorResponse{Message: "Error creating magic link"})
// 	}
// 	err = mailer.SendOnboardingMail(email, magicLinkToken)
// 	return err
// }

func AuthRoutes(app *fiber.App) {
	app.Get("/auth/check", func(c *fiber.Ctx) error {

		return c.SendStatus(200)
	})

	app.Post("/auth/onboard", func(c *fiber.Ctx) error {

		body := RegisterRequest{}
		err := c.BodyParser(&body)
		if err != nil {
			return c.Status(400).JSON(&ErrorResponse{Message: "Invalid request body"})
		}

		claims, ok := clerk.SessionClaimsFromContext(c.Context())
		if !ok {
			fmt.Println("Claims not found")
			return c.Status(401).JSON(&fiber.Map{"message": "Unauthorized"})
		}

		clerkUser, err := user.Get(context.Background(), claims.Subject)
		if err != nil {
			fmt.Println("CLerk User not found", err)
			return c.Status(401).JSON(&fiber.Map{"message": "Unauthorized"})
		}
		user := &models.User{}

		user.Email = clerkUser.EmailAddresses[0].EmailAddress
		user.Name = body.Name
		user.DisplayName = body.DisplayName
		user.ClerkId = clerkUser.ID

		err = createNewUser(c, user)
		if err != nil {
			return err
		}

		clerkOrg, err := organization.Create(context.Background(), &organization.CreateParams{
			CreatedBy: &clerkUser.ID,
			Name:      &body.OrgName,
		})

		if err != nil {
			return c.Status(500).JSON(&ErrorResponse{Message: "Error creating organization"})
		}

		org := &models.Organization{
			Name:    body.OrgName,
			ClerkId: clerkOrg.ID,
		}
		err = database.CreateOrganization(org)
		if err != nil {
			return c.Status(500).JSON(&ErrorResponse{Message: "Error creating user"})
		}
		userOrg := &models.UserOrganization{
			UserId:         user.ID,
			OrganizationId: org.ID,
		}
		err = database.CreateUserOrganization(userOrg)
		if err != nil {
			return c.Status(500).JSON(&ErrorResponse{Message: "Error creating user"})
		}

		return c.Status(200).JSON(&fiber.Map{"user_id": user.ID, "org_id": org.ClerkId})
	})

	// app.Post("/auth/login", func(c *fiber.Ctx) error {
	// 	body := LoginRequest{}
	// 	err:=c.BodyParser(&body)
	// 	if err!=nil{
	// 		fmt.Println(err)
	// 		return c.Status(400).JSON(&ErrorResponse{Message: "Invalid request body"})
	// 	}

	// 	// Check if the user exists in the database

	// 	err = sendMagicLink(c, body.Email)
	// 	if err!=nil{
	// 		return c.Status(500).JSON(&ErrorResponse{Message: "Error sending magic link"})
	// 	}

	// 	return c.SendStatus(200)
	// 	// token,err:=database.CreateSessionToken(user.ID,nil)

	// 	// if err!=nil{
	// 	// 	return c.Status(500).JSON(&ErrorResponse{Message: "Error creating session"})
	// 	// }
	// 	// SetSessionCookie(c, token)
	// 	// return c.JSON(&user)

	// })

	// app.Get("/auth/verify", func(c *fiber.Ctx) error {
	// 	token:=c.Query("token")
	// 	fmt.Printf("Vrerifying token = %s",token)
	// 	if token==""{
	// 		return c.Status(400).JSON(&ErrorResponse{Message: "Invalid token"})
	// 	}
	// 	parsedToken,err:=utils.ValidateMagicLinkToken(token)
	// 	if err!=nil{
	// 		fmt.Println(err)
	// 		return c.Status(400).JSON(&ErrorResponse{Message: "Invalid token"})
	// 	}

	// 	claims:=parsedToken.Claims.(jwt.MapClaims)
	// 	email:=claims["email"].(string)
	// 	fmt.Printf("Email = %s",email)

	// 	user:=&models.User{Email:email}

	// 	err = database.GetUserByEmail(email, user)

	// 	if err!=nil{
	// 		return createNewUser(c, user)
	// 	}

	// 	token,err=database.CreateSessionToken(user.ID, nil)
	// 	if err!=nil{
	// 		return c.Status(500).JSON(&ErrorResponse{Message: "Error creating session"})
	// 	}

	// 	SetSessionCookie(c, token)
	// 	if user.Name==""{
	// 		return c.Status(201).JSON(&fiber.Map{"message":"User created. Please complete registration"})
	// 	}
	// 	return c.Status(200).JSON(&user)
	// })

	// app.Post("/auth/register", func(c *fiber.Ctx) error {
	// 	token:=c.Cookies("audlib")
	// 	if token==""{
	// 		return c.Status(400).JSON(&ErrorResponse{Message: "token not provided"})
	// 	}

	// 	session:=&models.Session{}
	// 	err:= database.GetSessionByToken(token, session)

	// 	if err!=nil{
	// 		return c.Status(401).JSON(&fiber.Map{"message":"session with this token does not exist"})
	// 	}

	// 	if session.ExpiresAt<time.Now().Unix(){
	// 		c.ClearCookie("audlib")
	// 		return c.Status(401).JSON(fiber.Map{"message": "Session expired"})
	// 	}

	// 	body:=&RegisterRequest{}
	// 	err=c.BodyParser(body)

	// 	if err!=nil{
	// 		return c.Status(400).JSON(&ErrorResponse{Message: "Invalid request body"})
	// 	}

	// 	user:=session.User
	// 	user.Name=body.Name
	// 	user.DisplayName=body.DisplayName
	// 	config.DB.Save(user)

	// org:=&models.Organization{
	// 	Name: body.OrgName,
	// }
	// err= database.CreateOrganization(org)
	// if err!=nil{
	// 	return c.Status(500).JSON(&ErrorResponse{Message: "Error creating user"})
	// }
	// userOrg:=&models.UserOrganization{
	// 	UserId: session.User.ID,
	// 	OrganizationId: org.ID,
	// }
	// err=database.CreateUserOrganization(userOrg)
	// if err!=nil{
	// 	return c.Status(500).JSON(&ErrorResponse{Message: "Error creating user"})
	// }

	// 	err=database.ChangeOrganization(token,org.ID)
	// 	if err!=nil{
	// 		fmt.Println(err)
	// 		return c.Status(500).JSON(&ErrorResponse{Message: "cannot change organization"})
	// 	}

	// 	newToken,err:=database.CreateSessionToken(session.User.ID, nil)

	// 	if err!=nil{
	// 		return c.Status(500).JSON(&ErrorResponse{Message: "Error creating session"})
	// 	}

	// 	SetSessionCookie(c, newToken)

	// 	return c.JSON(&session.User)
	// })

	app.Get("/auth/info", func(c *fiber.Ctx) error {
		user, err := GetAuthenticatedUser(c)
		if err != nil {
			return c.Status(401).JSON(&ErrorResponse{Message: "No Authenticated user"})
		}
		return c.JSON(&user)
	})
	// app.Post("/auth/logout", func(c *fiber.Ctx) error {
	// 	c.ClearCookie("session")
	// 	return c.SendStatus(200)
	// })
}
