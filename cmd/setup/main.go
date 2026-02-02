package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"project-go/config"
)

func main() {
	log.Println("Setting up database...")
	
	// Load configuration
	config.LoadConfig()
	
	// First, connect to postgres database to create school database
	err := createDatabase()
	if err != nil {
		log.Printf("Warning: Could not create database (might already exist): %v", err)
	}
	
	// Now connect to school database
	config.InitDB()
	
	// Create tables manually using raw SQL
	err = createTables()
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
	
	// Insert initial data
	err = insertInitialData()
	if err != nil {
		log.Fatal("Failed to insert initial data:", err)
	}
	
	log.Println("Database setup completed successfully!")
}

func createDatabase() error {
	// Connect to default postgres database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBPort,
	)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	
	// Create school database
	return db.Exec("CREATE DATABASE school").Error
}

func createTables() error {
	// Create roles table
	err := config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`).Error
	if err != nil {
		return err
	}
	
	// Create permissions table
	err = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS permissions (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			description TEXT,
			resource VARCHAR(255) NOT NULL,
			action VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`).Error
	if err != nil {
		return err
	}
	
	// Create role_permissions junction table
	err = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS role_permissions (
			role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
			permission_id INTEGER REFERENCES permissions(id) ON DELETE CASCADE,
			PRIMARY KEY (role_id, permission_id)
		);
	`).Error
	if err != nil {
		return err
	}
	
	// Create users table
	err = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			role_id INTEGER NOT NULL REFERENCES roles(id),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL
		);
	`).Error
	if err != nil {
		return err
	}
	
	// Create indexes
	config.DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);")
	config.DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);")
	config.DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);")
	
	log.Println("Tables created successfully!")
	return nil
}

func insertInitialData() error {
	// Insert roles
	roles := []map[string]interface{}{
		{"name": "admin", "description": "Administrator with full access"},
		{"name": "teacher", "description": "Teacher with limited access"},
		{"name": "student", "description": "Student with basic access"},
	}
	
	for _, role := range roles {
		config.DB.Exec(`
			INSERT INTO roles (name, description) VALUES (?, ?) 
			ON CONFLICT (name) DO NOTHING
		`, role["name"], role["description"])
	}
	
	// Insert permissions
	permissions := []map[string]interface{}{
		{"name": "user.create", "description": "Create new users", "resource": "user", "action": "create"},
		{"name": "user.read", "description": "View users", "resource": "user", "action": "read"},
		{"name": "user.update", "description": "Update users", "resource": "user", "action": "update"},
		{"name": "user.delete", "description": "Delete users", "resource": "user", "action": "delete"},
		{"name": "role.create", "description": "Create new roles", "resource": "role", "action": "create"},
		{"name": "role.read", "description": "View roles", "resource": "role", "action": "read"},
		{"name": "role.update", "description": "Update roles", "resource": "role", "action": "update"},
		{"name": "role.delete", "description": "Delete roles", "resource": "role", "action": "delete"},
	}
	
	for _, perm := range permissions {
		config.DB.Exec(`
			INSERT INTO permissions (name, description, resource, action) VALUES (?, ?, ?, ?) 
			ON CONFLICT (name) DO NOTHING
		`, perm["name"], perm["description"], perm["resource"], perm["action"])
	}
	
	// Assign permissions to roles
	// Admin gets all permissions
	config.DB.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id FROM roles r, permissions p 
		WHERE r.name = 'admin'
		ON CONFLICT DO NOTHING
	`)
	
	// Teacher gets limited permissions
	config.DB.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id FROM roles r, permissions p 
		WHERE r.name = 'teacher' AND p.name IN ('user.read', 'role.read')
		ON CONFLICT DO NOTHING
	`)
	
	// Student gets basic permissions
	config.DB.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id FROM roles r, permissions p 
		WHERE r.name = 'student' AND p.name IN ('user.read')
		ON CONFLICT DO NOTHING
	`)
	
	log.Println("Initial data inserted successfully!")
	return nil
}