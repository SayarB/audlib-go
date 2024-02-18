package database

import (
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
)

func CreateProject(project *models.Project) error {
	return config.DB.Create(project).Error
}