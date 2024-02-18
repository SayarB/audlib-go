package models

import (
	"io"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	ID string `gorm:"primaryKey"`
}

type AudioFile struct {
	Model // This is the primary key
	BucketId  string 
	Folder    string 
	Key       string  `gorm:"primaryKey"`
	Extension string 
	File io.Reader `gorm:"-"`
	MIMEType  string
	ProjectId string
	Project   *Project
}

type Project struct{
	Model
	Name string
	OwnerId string
	Owner *Organization
	AudioFiles []AudioFile `gorm:"foreignKey:ProjectId;references:ID"`
}
type Organization struct{
	Model
	Name string
	Projects []Project `gorm:"foreignKey:OwnerId;references:ID"`
	
}

type User struct{
	Model
	Name string
	Email string
}

type UserOrganization struct{
	Model
	UserId string
	OrganizationId string
	Role string
}



func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	m.ID = id.String()
	return
}
