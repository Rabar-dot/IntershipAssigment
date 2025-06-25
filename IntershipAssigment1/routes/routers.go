package routes

import (
    "github.com/gin-gonic/gin"
    "golang-user-api/controllers"
    "golang-user-api/middleware"
)

func RegisterRoutes(r *gin.Engine) {
    api := r.Group("/api")
    api.POST("/login", controllers.Login)

    user := api.Group("/users")
    user.Use(middleware.AuthMiddleware())

    user.POST("/", middleware.AuthorizeRole("admin"), controllers.CreateUser)
    user.PUT("/:id", middleware.AuthorizeRole("admin"), controllers.UpdateUser)
    user.DELETE("/:id", middleware.AuthorizeRole("admin"), controllers.DeleteUser)
    user.POST("/:id/upload", middleware.AuthorizeRole("admin"), controllers.UploadImage)
    user.DELETE("/:id/image", middleware.AuthorizeRole("admin"), controllers.DeleteImage)

    user.GET("/", controllers.GetAllUsers)
    user.GET("/:id", controllers.GetUserByID)
}
