package database

import (
	"fmt"

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

func GetOrganizationsForUser(userId string) ([]models.UserOrganization, error){
	var orgs []models.UserOrganization
	tx:=config.DB.Where("user_id = ?", userId).Preload("Organization").Find(&orgs)
	fmt.Println(tx.RowsAffected)

	return orgs, tx.Error
}