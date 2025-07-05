// TEMPORARY: Allow user creation without authentication for initial setup
// package routes

// import (
// 	"golang-user-api/controllers"
// 	"golang-user-api/middleware"

// 	"github.com/gin-gonic/gin"
// )

// // RegisterRoutes sets up the API routes.
// func RegisterRoutes(r *gin.Engine) {
// 	api := r.Group("/api")                // Create API group
// 	api.POST("/login", controllers.Login) // Login route

// 	// TEMPORARY: Allow user creation without authentication for initial setup
// 	// This line moves the CreateUser route outside the authenticated group
// 	api.POST("/users", controllers.CreateUser) // <--- THIS IS THE KEY CHANGE

// 	user := api.Group("/users")           // User routes
// 	user.Use(middleware.AuthMiddleware()) // Apply authentication middleware
// 	// user.POST("/", controllers.CreateUser) // <--- COMMENT OUT OR REMOVE THIS LINE
// 	user.PUT("/:id", middleware.AuthorizeRole("admin"), controllers.UpdateUser)           // Update user route
// 	user.DELETE("/:id", middleware.AuthorizeRole("admin"), controllers.DeleteUser)        // Delete user route
// 	user.POST("/:id/upload", middleware.AuthorizeRole("admin"), controllers.UploadImage)  // Upload image route
// 	user.DELETE("/:id/image", middleware.AuthorizeRole("admin"), controllers.DeleteImage) // Delete image route
// 	user.GET("/", controllers.GetAllUsers)                                                // Get all users route
// 	user.GET("/:id", controllers.GetUserByID)                                             // Get user by ID route
// }

package routes

import (
	"golang-user-api/controllers"
	"golang-user-api/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the API routes.
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")                // Create API group
	api.POST("/login", controllers.Login) // Login route

	user := api.Group("/users")            // User routes
	user.Use(middleware.AuthMiddleware())  // Apply authentication middleware
	user.POST("/", controllers.CreateUser) // Create user route (now protected)
	user.PUT("/:id", middleware.AuthorizeRole("admin"), controllers.UpdateUser)
	user.DELETE("/:id", middleware.AuthorizeRole("admin"), controllers.DeleteUser)
	user.POST("/:id/upload", middleware.AuthorizeRole("admin"), controllers.UploadImage)
	user.DELETE("/:id/image", middleware.AuthorizeRole("admin"), controllers.DeleteImage)
	user.GET("/", controllers.GetAllUsers)    // Get all users route (protected)
	user.GET("/:id", controllers.GetUserByID) // Get user by ID route (protected)
}
