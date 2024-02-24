package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


var (
	streamTokenSecret=[]byte(os.Getenv("STREAM_TOKEN_SECRET"))
)

func GenerateStreamToken(key string ) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"key": key,
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

