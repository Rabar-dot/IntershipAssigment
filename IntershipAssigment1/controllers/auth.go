package controllers

import (
	"golang-user-api/config"
	"golang-user-api/models"
	"golang-user-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login handles user login and generates a JWT token.
func Login(c *gin.Context) {
	var body struct {
		Username, Password string
	}
	if c.BindJSON(&body) != nil {
		config.Logger.WithField("error", "Invalid input").Error("Failed to bind JSON for Login") // Add this
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	config.Logger.WithField("attempted_username", body.Username).Info("Login attempt for username") // Add this
	var u models.User
	if err := config.DB.Where("username = ?", body.Username).First(&u).Error; err != nil {
		config.Logger.WithField("username", body.Username).WithField("db_error", err.Error()).Error("User not found in DB during login") // Modify this line
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Use another Username"})
		return
	}
	if !utils.CheckPassword(u.Password, body.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Use another Passowrd"}) // Return error for password mismatch
		return
	}
	token, _ := utils.GenerateJWT(u.Username, u.Role)                                    // Generate JWT token
	c.JSON(http.StatusOK, gin.H{"token": token, "role": u.Role, "username": u.Username}) // Return token and user info
}
