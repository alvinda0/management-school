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
	}
}