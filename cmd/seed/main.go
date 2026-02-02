package main

import (
	"log"

	"project-go/config"
	"project-go/internal/migration"
	"project-go/internal/seeder"
)

func main() {
	log.Println("Starting database seeding...")
	
	// Load configuration
	config.LoadConfig()
	
	// Initialize database
	config.InitDB()
	
	// Run migrations first
	err := migration.RunMigrations()
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	
	// Run seeder
	seeder.SeedDatabase()
	
	log.Println("Seeding completed!")
}