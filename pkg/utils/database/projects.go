package database

import (
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
)

func CreateProject(project *models.Project) error {
	return config.DB.Create(project).Error
}

func GetProjectById(id string) (*models.Project, error) {
	project := &models.Project{}
	err := config.DB.Where("id = ?", id).First(project).Error
	return project, err
}
func GetProjectsByOrganizationId(id string) ([]models.Project, error) {
	var projects []models.Project
	err := config.DB.Where("owner_id = ?", id).Find(&projects).Error
	return projects, err
}