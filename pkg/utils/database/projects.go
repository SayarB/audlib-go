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
	err := config.DB.Where("id = ?", id).Preload("Versions").First(project).Error
	return project, err
}
func GetProjectsByOrganizationId(id string) ([]models.Project, error) {
	var projects []models.Project
	err := config.DB.Where("owner_id = ?", id).Find(&projects).Error
	return projects, err
}

func GetProjectsWithLatestVersion(id string) ([]models.ProjectWithLatestVersion, error) {
	var projectsWithLatestVersion []models.ProjectWithLatestVersion
	var projects []models.Project
	err := config.DB.Where("owner_id = ?", id).Preload("Versions").Find(&projects).Error
	for _, proj := range projects{
		if len(proj.Versions) == 0 {
			proj.Versions = nil
			projectsWithLatestVersion = append(projectsWithLatestVersion, models.ProjectWithLatestVersion{
				Project: proj,
				LatestVersion: nil,
			})
			continue
		}else{
			latestVersion:=proj.Versions[len(proj.Versions)-1]
			proj.Versions=nil
			projectsWithLatestVersion = append(projectsWithLatestVersion, models.ProjectWithLatestVersion{
				Project: proj,
				LatestVersion: &latestVersion,
			})
		}
	}
	return projectsWithLatestVersion, err
}