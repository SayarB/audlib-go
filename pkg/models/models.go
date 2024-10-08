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
	Model
	BucketId  string
	Folder    string
	Key       string `gorm:"unique"`
	Extension string
	File      io.Reader `gorm:"-"`
	Size      int64
	MIMEType  string
	AuthorId  string
	Author    *User
	Version   *Version `gorm:"foreignKey:AudioFileId;references:ID"`
}
type ProjectFile struct {
	Model
	BucketId  string
	Folder    string
	Key       string `gorm:"unique"`
	Extension string
	File      io.Reader `gorm:"-"`
	Size      int64
	MIMEType  string
	AuthorId  string
	Author    *User
	Version   *Version `gorm:"foreignKey:ProjectFileId;references:ID"`
}

type Version struct {
	Model
	Title         string
	AudioFileId   string
	ProjectFileId string
	// AudioFile *AudioFile
	ProjectId   string
	Project     *Project
	IsPublished bool
	AuthorId    string
	Author      *User
}

type Project struct {
	Model
	Name     string
	OwnerId  string
	Owner    *Organization
	Versions []Version `gorm:"foreignKey:ProjectId;references:ID"`
}
type Organization struct {
	Model
	Name     string
	Projects []Project          `gorm:"foreignKey:OwnerId;references:ID"`
	Users    []UserOrganization `gorm:"foreignKey:OrganizationId;references:ID"`
	Sessions []Session          `gorm:"foreignKey:OrganizationId;references:ID" `
	ClerkId  string             `gorm:"uniqueIndex"`
}

type User struct {
	Model
	Name        string
	ClerkId     string `gorm:"uniqueIndex"`
	DisplayName string `gorm:"uniqueIndex"`
	Email       string `gorm:"uniqueIndex"`
	// Password string
	Sessions      []Session          `gorm:"foreignKey:UserId;references:ID"`
	Organizations []UserOrganization `gorm:"foreignKey:UserId;references:ID"`
	AudioFiles    []AudioFile        `gorm:"foreignKey:AuthorId;references:ID"`
	Versions      []Version          `gorm:"foreignKey:AuthorId;references:ID"`
	MagicLinks    []MagicLink        `gorm:"foreignKey:UserId;references:ID"`
}

type Session struct {
	Model
	UserId         string
	Token          string
	ExpiresAt      int64
	OrganizationId *string
	Organization   *Organization
	User           *User
}

type UserOrganization struct {
	Model
	UserId         string
	User           *User
	OrganizationId string
	Organization   *Organization
	Role           string
}

type MagicLink struct {
	Token  string
	UserId string
	User   *User
	Used   bool
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	m.ID = id.String()
	return
}

type ProjectWithLatestVersion struct {
	Project
	LatestVersion *Version
}
