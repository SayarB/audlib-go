package database

import (
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
)

func GetOrganizationById(orgId string, org *models.Organization) error {
	return config.DB.Where("id = ?", orgId).First(org).Error
}

func CreateOrganization(org *models.Organization) error {
	return config.DB.Create(org).Error
}

func CreateUserOrganization(userOrg *models.UserOrganization) error {
	return config.DB.Create(userOrg).Error
}

func GetUserOrganization(userOrg *models.UserOrganization) error {
	return config.DB.Where("user_id = ? AND organization_id = ?", userOrg.UserId, userOrg.OrganizationId).First(userOrg).Error
}