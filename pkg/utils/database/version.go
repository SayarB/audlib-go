package database

import (
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
)

func CreateVersion(version *models.Version) error {
	return config.DB.Create(version).Error
}


func GetVersionsByProjectId(projectId string) ([]models.Version, error) {
	var versions []models.Version
	err:=config.DB.Where("project_id = ?", projectId).Find(&versions).Error
	return versions,err
}

func GetVersionById(versionId string) (*models.Version, error) {
	var version models.Version
	err:=config.DB.Where("id = ?", versionId).First(&version).Error
	return &version,err
}