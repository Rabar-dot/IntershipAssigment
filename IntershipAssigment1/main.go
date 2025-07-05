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
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    config.InitLogger() // This is correctly placed
    config.InitDB()
    config.InitRedis()

    router := gin.Default()
    router.Static("/uploads", "./uploads")
    routes.RegisterRoutes(router)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Println("Server running on port", port)
    router.Run(":" + port)
}

