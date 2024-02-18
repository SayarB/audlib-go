package database

import (
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
	"gorm.io/gorm"
)

func CreateAudioFile(db *gorm.DB, audio *models.AudioFile) error {
	err := db.Create(audio).Error
	if err != nil {
		return err
	}
	return nil
}

func GetFileInfo(key string) (*models.AudioFile,error){
	fileInfo:=&models.AudioFile{}
	db:=config.DB.Where("key = ?", key).First(fileInfo)
	if db.Error!=nil{
		return nil, db.Error
	}
	return fileInfo,nil
}