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
	
	// Seed permissions
	seedPermissions()
	
	// Assign permissions to roles
	assignPermissionsToRoles()

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

func seedPermissions() {
	permissions := []model.Permission{
		// User permissions
		{Name: "user.create", Description: "Create new users", Resource: "user", Action: "create"},
		{Name: "user.read", Description: "View users", Resource: "user", Action: "read"},
		{Name: "user.update", Description: "Update users", Resource: "user", Action: "update"},
		{Name: "user.delete", Description: "Delete users", Resource: "user", Action: "delete"},
		
		// Role permissions
		{Name: "role.create", Description: "Create new roles", Resource: "role", Action: "create"},
		{Name: "role.read", Description: "View roles", Resource: "role", Action: "read"},
		{Name: "role.update", Description: "Update roles", Resource: "role", Action: "update"},
		{Name: "role.delete", Description: "Delete roles", Resource: "role", Action: "delete"},
		
		// Academic permissions
		{Name: "academic.manage", Description: "Manage academic data", Resource: "academic", Action: "manage"},
		{Name: "academic.read", Description: "View academic data", Resource: "academic", Action: "read"},
		
		// Student permissions
		{Name: "student.manage", Description: "Manage student data", Resource: "student", Action: "manage"},
		{Name: "student.read", Description: "View student data", Resource: "student", Action: "read"},
		
		// Teacher permissions
		{Name: "teacher.manage", Description: "Manage teacher data", Resource: "teacher", Action: "manage"},
		{Name: "teacher.read", Description: "View teacher data", Resource: "teacher", Action: "read"},
		
		// Class permissions
		{Name: "class.manage", Description: "Manage class data", Resource: "class", Action: "manage"},
		{Name: "class.read", Description: "View class data", Resource: "class", Action: "read"},
		
		// Report permissions
		{Name: "report.generate", Description: "Generate reports", Resource: "report", Action: "generate"},
		{Name: "report.read", Description: "View reports", Resource: "report", Action: "read"},
		
		// System permissions
		{Name: "system.manage", Description: "Manage system settings", Resource: "system", Action: "manage"},
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
	var systemRole, kepalaSekolahRole, staffRole, guruRole, siswaRole model.Role
	config.DB.Where("name = ?", "system").First(&systemRole)
	config.DB.Where("name = ?", "kepala_sekolah").First(&kepalaSekolahRole)
	config.DB.Where("name = ?", "staff").First(&staffRole)
	config.DB.Where("name = ?", "guru").First(&guruRole)
	config.DB.Where("name = ?", "siswa").First(&siswaRole)

	// Get all permissions
	var allPermissions []model.Permission
	config.DB.Find(&allPermissions)

	// System gets all permissions
	if systemRole.ID != 0 {
		config.DB.Model(&systemRole).Association("Permissions").Replace(allPermissions)
		log.Println("Assigned all permissions to system role")
	}

	// Kepala Sekolah gets management permissions
	if kepalaSekolahRole.ID != 0 {
		var kepalaSekolahPermissions []model.Permission
		config.DB.Where("name IN ?", []string{
			"user.create", "user.read", "user.update", "user.delete",
			"role.read", "academic.manage", "academic.read",
			"student.manage", "student.read", "teacher.manage", "teacher.read",
			"class.manage", "class.read", "report.generate", "report.read",
		}).Find(&kepalaSekolahPermissions)
		config.DB.Model(&kepalaSekolahRole).Association("Permissions").Replace(kepalaSekolahPermissions)
		log.Println("Assigned management permissions to kepala sekolah role")
	}

	// Staff gets administrative permissions
	if staffRole.ID != 0 {
		var staffPermissions []model.Permission
		config.DB.Where("name IN ?", []string{
			"user.read", "user.update", "academic.read",
			"student.manage", "student.read", "teacher.read",
			"class.read", "report.read",
		}).Find(&staffPermissions)
		config.DB.Model(&staffRole).Association("Permissions").Replace(staffPermissions)
		log.Println("Assigned administrative permissions to staff role")
	}

	// Guru gets teaching permissions
	if guruRole.ID != 0 {
		var guruPermissions []model.Permission
		config.DB.Where("name IN ?", []string{
			"user.read", "academic.read", "student.read",
			"class.manage", "class.read", "report.read",
		}).Find(&guruPermissions)
		config.DB.Model(&guruRole).Association("Permissions").Replace(guruPermissions)
		log.Println("Assigned teaching permissions to guru role")
	}

	// Siswa gets basic permissions
	if siswaRole.ID != 0 {
		var siswaPermissions []model.Permission
		config.DB.Where("name IN ?", []string{
			"academic.read", "class.read", "report.read",
		}).Find(&siswaPermissions)
		config.DB.Model(&siswaRole).Association("Permissions").Replace(siswaPermissions)
		log.Println("Assigned basic permissions to siswa role")
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
	// Import crypto/bcrypt at the top of the file
	const cost = 12
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}