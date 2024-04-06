package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sayar/go-streaming/pkg/models"
	"github.com/sayar/go-streaming/pkg/utils/database"
)


var (
	streamTokenSecret=[]byte(os.Getenv("STREAM_TOKEN_SECRET"))
)

func GenerateStreamToken(version *models.Version ) (string, error) {
	fileInfo, err := database.GetFileInfo(version.AudioFileId)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"vid": version.ID,
		"key": fileInfo.Key,
		"mime": fileInfo.MIMEType,
		"size":fileInfo.Size,
		"bucket":fileInfo.BucketId,
		"extension":fileInfo.Extension,
		"exp":time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString(streamTokenSecret)
}


func ValidateStreamToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return streamTokenSecret, nil
	})
	return parsedToken, err
}

func GenerateMagicLinkToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"purpose": "magic_link",
		"exp": time.Now().Add(5*time.Minute).Unix(),
	})
	return token.SignedString(streamTokenSecret)
}

func ValidateMagicLinkToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return streamTokenSecret, nil
	})
	return parsedToken, err
}