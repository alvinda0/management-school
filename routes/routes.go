package routes

import (
	"github.com/gin-gonic/gin"
	"project-go/internal/handler"
	"project-go/internal/repository"
	"project-go/internal/service"
	"project-go/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// Initialize dependencies
	userRepo := repository.NewUserRepository()
	roleRepo := repository.NewRoleRepository()
	userService := service.NewUserService(userRepo, roleRepo)
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleRepo)

	// Public routes
	api := r.Group("/api/v1")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// User management routes - basic access for all authenticated users
		protected.GET("/users", userHandler.GetAllUsers)
		protected.GET("/users/:id", userHandler.GetUser)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)
		
		// Role management routes - basic access for all authenticated users
		protected.GET("/roles", roleHandler.GetAllRoles)
		protected.GET("/roles/:id", roleHandler.GetRole)
		
		// Admin-only routes (system and kepala sekolah only)
		adminOnly := protected.Group("/admin")
		adminOnly.Use(middleware.IsSystemOrKepalaSekolah())
		{
			adminOnly.POST("/users", userHandler.Register) // Create users
			adminOnly.POST("/roles", roleHandler.CreateRole) // Create roles
			adminOnly.PUT("/roles/:id", roleHandler.UpdateRole) // Update roles
			adminOnly.DELETE("/roles/:id", roleHandler.DeleteRole) // Delete roles
		}
		
		// Management routes (kepala sekolah, staff, guru can access)
		management := protected.Group("/management")
		management.Use(middleware.RequireAnyRole("kepala_sekolah", "staff", "guru"))
		{
			// Academic management
			management.GET("/academic", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Academic data"})
			})
			management.POST("/academic", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Academic data created"})
			})
			
			// Student management
			management.GET("/students", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Student data"})
			})
			management.POST("/students", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Student created"})
			})
		}
		
		// Teacher-only routes
		teacherOnly := protected.Group("/teacher")
		teacherOnly.Use(middleware.RequireRole("guru"))
		{
			teacherOnly.GET("/classes", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Class data for teacher"})
			})
			teacherOnly.POST("/classes", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Class created by teacher"})
			})
		}
		
		// Student-only routes
		studentOnly := protected.Group("/student")
		studentOnly.Use(middleware.RequireRole("siswa"))
		{
			studentOnly.GET("/profile", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Student profile"})
			})
			studentOnly.GET("/grades", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Student grades"})
			})
		}
	}
}