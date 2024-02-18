package database

import (
	"github.com/sayar/go-streaming/pkg/config"
	"github.com/sayar/go-streaming/pkg/models"
)

func GetOrganizationById(orgId string, org *models.Organization) error {
	return config.DB.Where("id = ?", orgId).First(org).Error
}