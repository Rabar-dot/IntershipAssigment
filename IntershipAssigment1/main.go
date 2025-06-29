package main

import (
    "log"
    "os"

    "github.com/joho/godotenv"
    "github.com/gin-gonic/gin"

    "golang-user-api/config"
    "golang-user-api/routes"
)

func main() {
    err := godotenv.Load() // Load environment variables from .env file
    if err != nil {
        log.Fatal("Error loading .env file") // Log error if .env file cannot be loaded
    }

    config.InitLogger() // Initialize logger
    config.InitDB() // Initialize database connection
    config.InitRedis() // Initialize Redis connection

    router := gin.Default() // Create a new Gin router
    router.Static("/uploads", "./uploads") // Serve static files from uploads directory
    routes.RegisterRoutes(router) // Register API routes

    port := os.Getenv("PORT") // Get port from environment variable
    if port == "" {
        port = "8080" // Default to port 8080 if not specified
    }
    log.Println("Server running on port", port) // Log server start
    router.Run(":" + port) // Start the server
}
