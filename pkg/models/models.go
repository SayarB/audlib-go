package models

import (
	"io"

	"gorm.io/gorm"
)

type AudioFile struct {
	gorm.Model // This is the primary key
	BucketId  string 
	Folder    string 
	Key       string  `gorm:"primaryKey"`
	Extension string 
	File io.Reader `gorm:"-"`
	MIMEType  string
}

type ErrorResponse struct {
	Message string `json:"message"`
}