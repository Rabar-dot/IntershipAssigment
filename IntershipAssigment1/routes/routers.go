package routes

import (
    "github.com/gin-gonic/gin"
    "golang-user-api/controllers"
    "golang-user-api/middleware"
)

// RegisterRoutes sets up the API routes.
func RegisterRoutes(r *gin.Engine) {
    api := r.Group("/api") // Create API group
    api.POST("/login", controllers.Login) // Login route

    user := api.Group("/users") // User routes
    user.Use(middleware.AuthMiddleware()) // Apply authentication middleware

    user.POST("/", middleware.AuthorizeRole("admin"), controllers.CreateUser) // Create user route
    user.PUT("/:id", middleware.AuthorizeRole("admin"), controllers.UpdateUser) // Update user route
    user.DELETE("/:id", middleware.AuthorizeRole("admin"), controllers.DeleteUser) // Delete user route
    user.POST("/:id/upload", middleware.AuthorizeRole("admin"), controllers.UploadImage) // Upload image route
    user.DELETE("/:id/image", middleware.AuthorizeRole("admin"), controllers.DeleteImage) // Delete image route

    user.GET("/", controllers.GetAllUsers) // Get all users route
    user.GET("/:id", controllers.GetUserByID) // Get user by ID route
}
