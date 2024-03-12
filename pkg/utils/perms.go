package utils

import (
	"github.com/sayar/go-streaming/pkg/models"
)

func HasPermissionToStream(fileInfo *models.Version) bool {
	return true
}