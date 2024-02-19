package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
)

func CreateSession(session *models.Session) error {
	return config.DB.Create(session).Error
}

func CreateSessionToken(userId string) (string, error) {
	token,_:=uuid.NewV7()

	session := &models.Session{
		UserId: userId,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 3).Unix(),
		Token: token.String(),
	}
	err:=CreateSession(session)
	return token.String(), err
}