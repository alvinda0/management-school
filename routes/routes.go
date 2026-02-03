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
		// User management routes with permission checks
		protected.GET("/users", middleware.RequirePermission("user.read"), userHandler.GetAllUsers)
		protected.GET("/users/:id", middleware.RequirePermission("user.read"), userHandler.GetUser)
		protected.PUT("/users/:id", middleware.RequirePermission("user.update"), userHandler.UpdateUser)
		protected.DELETE("/users/:id", middleware.RequirePermission("user.delete"), userHandler.DeleteUser)
		
		// Role management routes
		protected.GET("/roles", middleware.RequirePermission("role.read"), roleHandler.GetAllRoles)
		protected.GET("/roles/:id", middleware.RequirePermission("role.read"), roleHandler.GetRole)
		
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
			management.GET("/academic", middleware.RequirePermission("academic.read"), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Academic data"})
			})
			management.POST("/academic", middleware.RequirePermission("academic.manage"), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Academic data created"})
			})
			
			// Student management
			management.GET("/students", middleware.RequirePermission("student.read"), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Student data"})
			})
			management.POST("/students", middleware.RequirePermission("student.manage"), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Student created"})
			})
		}
		
		// Teacher-only routes
		teacherOnly := protected.Group("/teacher")
		teacherOnly.Use(middleware.RequireRole("guru"))
		{
			teacherOnly.GET("/classes", middleware.RequirePermission("class.read"), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Class data for teacher"})
			})
			teacherOnly.POST("/classes", middleware.RequirePermission("class.manage"), func(c *gin.Context) {
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
			studentOnly.GET("/grades", middleware.RequirePermission("academic.read"), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Student grades"})
			})
		}
	}
}