package database

import (
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
)

func CreateAudioFile(audio *models.AudioFile) error {
	err := config.DB.Create(audio).Error
	if err != nil {
		return err
	}
	return nil
}
func CreateProjectFile(projectFile *models.ProjectFile) error {
	err := config.DB.Create(projectFile).Error
	if err != nil {
		return err
	}
	return nil
}

func GetFileInfo(id string) (*models.AudioFile, error) {
	fileInfo := &models.AudioFile{}
	db := config.DB.Where("id = ?", id).First(fileInfo)
	if db.Error != nil {
		return nil, db.Error
	}
	return fileInfo, nil
}
