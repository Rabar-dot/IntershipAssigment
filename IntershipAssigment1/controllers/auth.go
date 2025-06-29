package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "golang-user-api/config"
    "golang-user-api/models"
    "golang-user-api/utils"
)

// Login handles user login and generates a JWT token.
func Login(c *gin.Context) {
    var body struct{ Username, Password string }
    if c.BindJSON(&body) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // Return error for invalid input
        return
    }
    var u models.User
    if err := config.DB.Where("username = ?", body.Username).First(&u).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) // Return error for invalid credentials
        return
    }
    if !utils.CheckPassword(u.Password, body.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) // Return error for password mismatch
        return
    }
    token, _ := utils.GenerateJWT(u.Username, u.Role) // Generate JWT token
    c.JSON(http.StatusOK, gin.H{"token": token, "role": u.Role, "username": u.Username}) // Return token and user info
}
