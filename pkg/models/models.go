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
	Users []UserOrganization `gorm:"foreignKey:OrganizationId;references:ID"`
	
}

type User struct{
	Model
	Name string
	Email string `gorm:"unique"`
	Password string
	Sessions []Session `gorm:"foreignKey:UserId;references:ID"`
	Organizations []UserOrganization `gorm:"foreignKey:UserId;references:ID"`
}

type Session struct{
	Model
	UserId string
	Token string
	ExpiresAt int64
	User *User
}

type UserOrganization struct{
	Model
	UserId string
	User *User
	OrganizationId string
	Organization *Organization
	Role string
}



func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	m.ID = id.String()
	return
}
