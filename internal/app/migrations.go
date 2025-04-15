package app

import (
	"log"

	"trello-backend/internal/models"

	"gorm.io/gorm"
)

// Migrate runs the database migrations
func Migrate(db *gorm.DB) {
	log.Println("Running migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Board{},
		&models.List{},
		&models.Card{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully.")
}
