package seeder

import (
	"log"

	"project-go/config"
	"project-go/internal/model"
)

func SeedDatabase() {
	log.Println("Starting database seeding...")

	// Seed roles
	seedRoles()
	
	// Seed permissions
	seedPermissions()
	
	// Assign permissions to roles
	assignPermissionsToRoles()

	log.Println("Database seeding completed!")
}

func seedRoles() {
	roles := []model.Role{
		{
			Name:        "admin",
			Description: "Administrator with full access",
		},
		{
			Name:        "teacher",
			Description: "Teacher with limited access",
		},
		{
			Name:        "student",
			Description: "Student with basic access",
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

func seedPermissions() {
	permissions := []model.Permission{
		{Name: "user.create", Description: "Create new users", Resource: "user", Action: "create"},
		{Name: "user.read", Description: "View users", Resource: "user", Action: "read"},
		{Name: "user.update", Description: "Update users", Resource: "user", Action: "update"},
		{Name: "user.delete", Description: "Delete users", Resource: "user", Action: "delete"},
		{Name: "role.create", Description: "Create new roles", Resource: "role", Action: "create"},
		{Name: "role.read", Description: "View roles", Resource: "role", Action: "read"},
		{Name: "role.update", Description: "Update roles", Resource: "role", Action: "update"},
		{Name: "role.delete", Description: "Delete roles", Resource: "role", Action: "delete"},
	}

	for _, permission := range permissions {
		var existingPermission model.Permission
		result := config.DB.Where("name = ?", permission.Name).First(&existingPermission)
		if result.Error != nil {
			// Permission doesn't exist, create it
			if err := config.DB.Create(&permission).Error; err != nil {
				log.Printf("Error creating permission %s: %v", permission.Name, err)
			} else {
				log.Printf("Created permission: %s", permission.Name)
			}
		} else {
			log.Printf("Permission %s already exists", permission.Name)
		}
	}
}

func assignPermissionsToRoles() {
	// Get roles
	var adminRole, teacherRole, studentRole model.Role
	config.DB.Where("name = ?", "admin").First(&adminRole)
	config.DB.Where("name = ?", "teacher").First(&teacherRole)
	config.DB.Where("name = ?", "student").First(&studentRole)

	// Get all permissions
	var allPermissions []model.Permission
	config.DB.Find(&allPermissions)

	// Admin gets all permissions
	if adminRole.ID != 0 {
		config.DB.Model(&adminRole).Association("Permissions").Replace(allPermissions)
		log.Println("Assigned all permissions to admin role")
	}

	// Teacher gets limited permissions
	if teacherRole.ID != 0 {
		var teacherPermissions []model.Permission
		config.DB.Where("name IN ?", []string{"user.read", "role.read"}).Find(&teacherPermissions)
		config.DB.Model(&teacherRole).Association("Permissions").Replace(teacherPermissions)
		log.Println("Assigned limited permissions to teacher role")
	}

	// Student gets basic permissions
	if studentRole.ID != 0 {
		var studentPermissions []model.Permission
		config.DB.Where("name IN ?", []string{"user.read"}).Find(&studentPermissions)
		config.DB.Model(&studentRole).Association("Permissions").Replace(studentPermissions)
		log.Println("Assigned basic permissions to student role")
	}
}