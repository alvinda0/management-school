package seeder

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"project-go/config"
	"project-go/internal/model"
)

func SeedDatabase() {
	log.Println("Starting database seeding...")

	// Seed roles
	seedRoles()

	// Seed default users
	seedUsers()

	log.Println("Database seeding completed!")
}

func seedRoles() {
	roles := []model.Role{
		{
			Name:        "system",
			Description: "System administrator with full access",
		},
		{
			Name:        "kepala_sekolah",
			Description: "Kepala sekolah dengan akses manajemen penuh",
		},
		{
			Name:        "staff",
			Description: "Staff administrasi dengan akses terbatas",
		},
		{
			Name:        "guru",
			Description: "Guru dengan akses pembelajaran dan siswa",
		},
		{
			Name:        "siswa",
			Description: "Siswa dengan akses dasar",
		},
	}

	for _, role := range roles {
		var existingRole model.Role
		result := config.DB.Where("name = ?", role.Name).First(&existingRole)
		if result.Error != nil {
			// Role doesn't exist, create it
			if err := config.DB.Create(&role).Error; err != nil {
				log.Printf("Error creating role %s: %v", role.Name, err)
			} else {
				log.Printf("Created role: %s", role.Name)
			}
		} else {
			log.Printf("Role %s already exists", role.Name)
		}
	}
}

func seedUsers() {
	// Get system role
	var systemRole model.Role
	result := config.DB.Where("name = ?", "system").First(&systemRole)
	if result.Error != nil {
		log.Printf("System role not found: %v", result.Error)
		return
	}

	// Check if system user already exists
	var existingUser model.User
	result = config.DB.Where("email = ?", "system@gmail.com").First(&existingUser)
	if result.Error == nil {
		log.Println("System user already exists")
		return
	}

	// Hash password
	hashedPassword, err := hashPassword("system123")
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return
	}

	// Create system user
	systemUser := model.User{
		Username: "system",
		Email:    "system@gmail.com",
		Password: hashedPassword,
		RoleID:   systemRole.ID,
	}

	if err := config.DB.Create(&systemUser).Error; err != nil {
		log.Printf("Error creating system user: %v", err)
	} else {
		log.Println("Created system user successfully")
	}
}

func hashPassword(password string) (string, error) {
	const cost = 12
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}