package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
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
		{"name": "system", "description": "System administrator with full access"},
		{"name": "kepala_sekolah", "description": "Kepala sekolah dengan akses manajemen penuh"},
		{"name": "staff", "description": "Staff administrasi dengan akses terbatas"},
		{"name": "guru", "description": "Guru dengan akses pembelajaran dan siswa"},
		{"name": "siswa", "description": "Siswa dengan akses dasar"},
	}
	
	for _, role := range roles {
		config.DB.Exec(`
			INSERT INTO roles (name, description) VALUES (?, ?) 
			ON CONFLICT (name) DO NOTHING
		`, role["name"], role["description"])
	}
	
	// Insert permissions
	permissions := []map[string]interface{}{
		// User permissions
		{"name": "user.create", "description": "Create new users", "resource": "user", "action": "create"},
		{"name": "user.read", "description": "View users", "resource": "user", "action": "read"},
		{"name": "user.update", "description": "Update users", "resource": "user", "action": "update"},
		{"name": "user.delete", "description": "Delete users", "resource": "user", "action": "delete"},
		
		// Role permissions
		{"name": "role.create", "description": "Create new roles", "resource": "role", "action": "create"},
		{"name": "role.read", "description": "View roles", "resource": "role", "action": "read"},
		{"name": "role.update", "description": "Update roles", "resource": "role", "action": "update"},
		{"name": "role.delete", "description": "Delete roles", "resource": "role", "action": "delete"},
		
		// Academic permissions
		{"name": "academic.manage", "description": "Manage academic data", "resource": "academic", "action": "manage"},
		{"name": "academic.read", "description": "View academic data", "resource": "academic", "action": "read"},
		
		// Student permissions
		{"name": "student.manage", "description": "Manage student data", "resource": "student", "action": "manage"},
		{"name": "student.read", "description": "View student data", "resource": "student", "action": "read"},
		
		// Teacher permissions
		{"name": "teacher.manage", "description": "Manage teacher data", "resource": "teacher", "action": "manage"},
		{"name": "teacher.read", "description": "View teacher data", "resource": "teacher", "action": "read"},
		
		// Class permissions
		{"name": "class.manage", "description": "Manage class data", "resource": "class", "action": "manage"},
		{"name": "class.read", "description": "View class data", "resource": "class", "action": "read"},
		
		// Report permissions
		{"name": "report.generate", "description": "Generate reports", "resource": "report", "action": "generate"},
		{"name": "report.read", "description": "View reports", "resource": "report", "action": "read"},
		
		// System permissions
		{"name": "system.manage", "description": "Manage system settings", "resource": "system", "action": "manage"},
	}
	
	for _, perm := range permissions {
		config.DB.Exec(`
			INSERT INTO permissions (name, description, resource, action) VALUES (?, ?, ?, ?) 
			ON CONFLICT (name) DO NOTHING
		`, perm["name"], perm["description"], perm["resource"], perm["action"])
	}
	
	// Assign permissions to roles
	// System gets all permissions
	config.DB.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id FROM roles r, permissions p 
		WHERE r.name = 'system'
		ON CONFLICT DO NOTHING
	`)
	
	// Kepala Sekolah gets management permissions
	config.DB.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id FROM roles r, permissions p 
		WHERE r.name = 'kepala_sekolah' AND p.name IN (
			'user.create', 'user.read', 'user.update', 'user.delete',
			'role.read', 'academic.manage', 'academic.read',
			'student.manage', 'student.read', 'teacher.manage', 'teacher.read',
			'class.manage', 'class.read', 'report.generate', 'report.read'
		)
		ON CONFLICT DO NOTHING
	`)
	
	// Staff gets administrative permissions
	config.DB.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id FROM roles r, permissions p 
		WHERE r.name = 'staff' AND p.name IN (
			'user.read', 'user.update', 'academic.read',
			'student.manage', 'student.read', 'teacher.read',
			'class.read', 'report.read'
		)
		ON CONFLICT DO NOTHING
	`)
	
	// Guru gets teaching permissions
	config.DB.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id FROM roles r, permissions p 
		WHERE r.name = 'guru' AND p.name IN (
			'user.read', 'academic.read', 'student.read',
			'class.manage', 'class.read', 'report.read'
		)
		ON CONFLICT DO NOTHING
	`)
	
	// Siswa gets basic permissions
	config.DB.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id FROM roles r, permissions p 
		WHERE r.name = 'siswa' AND p.name IN (
			'academic.read', 'class.read', 'report.read'
		)
		ON CONFLICT DO NOTHING
	`)

	// Insert default system user
	err := insertSystemUser()
	if err != nil {
		log.Printf("Warning: Could not create system user: %v", err)
	}
	
	log.Println("Initial data inserted successfully!")
	return nil
}

func insertSystemUser() error {
	// Check if system user already exists
	var count int64
	config.DB.Raw("SELECT COUNT(*) FROM users WHERE email = ?", "system@gmail.com").Scan(&count)
	if count > 0 {
		log.Println("System user already exists")
		return nil
	}

	// Hash password
	hashedPassword, err := hashPassword("system123")
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	// Get system role ID
	var roleID uint
	err = config.DB.Raw("SELECT id FROM roles WHERE name = ?", "system").Scan(&roleID).Error
	if err != nil {
		return fmt.Errorf("system role not found: %v", err)
	}

	// Insert system user
	err = config.DB.Exec(`
		INSERT INTO users (username, email, password, role_id) 
		VALUES (?, ?, ?, ?)
	`, "system", "system@gmail.com", hashedPassword, roleID).Error
	
	if err != nil {
		return fmt.Errorf("error creating system user: %v", err)
	}

	log.Println("Created system user successfully")
	return nil
}

func hashPassword(password string) (string, error) {
	const cost = 12
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}