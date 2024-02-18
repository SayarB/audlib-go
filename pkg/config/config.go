package config

import (
	"fmt"
	"os"

	"github.com/sayar/go-streaming/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
func ConnectToDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=postgres.mpozbbvhwvezhtdhdeio password=%s host=aws-0-ap-south-1.pooler.supabase.com port=5432 dbname=postgres", os.Getenv("SUPABASE_DB_PASSWORD"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB=db
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to database")
	db.AutoMigrate(&models.AudioFile{})
	fmt.Println("Migrated")
	return db, nil
}
