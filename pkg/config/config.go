package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sayar/go-streaming/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var S3Client *s3.Client
func ConnectToDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=postgres.mpozbbvhwvezhtdhdeio password=%s host=aws-0-ap-south-1.pooler.supabase.com port=5432 dbname=postgres", os.Getenv("SUPABASE_DB_PASSWORD"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB=db
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to database")
	db.AutoMigrate(&models.AudioFile{}, &models.Organization{},&models.User{}, &models.Project{}, &models.Session{}, models.UserOrganization{}, &models.Version{})
	fmt.Println("Migrated")
	return db, nil
}


func ConfigureS3Client() error {
	s3Config,err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	S3Client=s3.NewFromConfig(s3Config)
	return err
}