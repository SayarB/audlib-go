package database

import (
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
)

func GetUserById(userId string, user *models.User) error {
	return config.DB.Where("id = ?", userId).First(user).Error
}

func GetUserByEmail(email string, user *models.User) error {
	return config.DB.Preload("Organizations.Organization").Where("email = ?", email).First(user).Error
}

func CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}