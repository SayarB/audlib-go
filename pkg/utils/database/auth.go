package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
)

func GetSessionByToken(token string, session *models.Session) error {
	return config.DB.Where("token = ?",token).Preload("User").Preload("Organization").First(session).Error
}

func CreateSession(session *models.Session) error {
	return config.DB.Create(session).Error
}

func UpdateSession(session *models.Session) error{
	newSession := models.Session{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 3).Unix(),
	}
	return config.DB.Where("id = ?", session.ID).Updates(newSession).Error
}

func CreateSessionToken(userId string, orgId *string) (string, error) {

	fmt.Println("creating new session token")

	token,_:=uuid.NewV7()
	
	newSession := &models.Session{
		UserId: userId,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 3).Unix(),
		Token: token.String(),
	}

	err:=CreateSession(newSession)
	return token.String(), err
}

func UpdateSessionToken(session *models.Session) (string, error) {

	err:=UpdateSession(session)
	return session.Token, err	
}