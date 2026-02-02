package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"project-go/config"
	"project-go/internal/migration"
	"project-go/internal/seeder"
	"project-go/routes"
)

func main() {
	// Load configuration
	config.LoadConfig()
	
	// Initialize database
	config.InitDB()
	
	// Run migrations
	err := migration.RunMigrations()
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	
	// Seed database with initial data
	seeder.SeedDatabase()
	
	// Setup router
	r := gin.Default()
	routes.SetupRoutes(r)
	
	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}