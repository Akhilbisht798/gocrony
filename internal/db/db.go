package db

import (
	"log"
	"os"

	"github.com/akhilbisht798/gocrony/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Job{})
	db.AutoMigrate(&models.Logs{})
	db.AutoMigrate(&models.UserIdentity{})

	DB = db
	log.Println("successfully connected to database!")
}
