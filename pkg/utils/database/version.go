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
func GetVersionByIdWithProject(versionId string) (*models.Version, error) {
	var version models.Version
	err:=config.DB.Preload("Project").Preload("Author").Where("id = ?", versionId).First(&version).Error
	return &version,err
}

func PublishVersion(version *models.Version) error {
	version.IsPublished=true
	return config.DB.Save(version).Error
}

func GetPublishedVersionByProjectId(projectId string) (*models.Version, error) {
	var version models.Version
	err:=config.DB.Where("project_id = ? AND is_published = ?", projectId, true).First(&version).Error
	return &version,err
}