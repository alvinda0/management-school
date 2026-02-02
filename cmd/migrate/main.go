package main

import (
	"log"

	"project-go/config"
	"project-go/internal/migration"
)

func main() {
	log.Println("Starting database migration...")
	
	// Load configuration
	config.LoadConfig()
	
	// Initialize database
	config.InitDB()
	
	// Run migrations
	err := migration.RunMigrations()
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	
	log.Println("Migration completed successfully!")
}