package controllers

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "golang-user-api/config"
    "golang-user-api/models"
    "golang-user-api/utils"
)

// CreateUser creates a new user in the database.
func CreateUser(c *gin.Context) {
    var u models.User
    if c.BindJSON(&u) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // Return error for invalid input
        return
    }
    if len(u.Username) == 0 || !isAlpha(u.Username) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"}) // Return error for invalid username
        return
    }
    if len(u.Password) < 8 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Password too short"}) // Return error for short password
        return
    }
    h, _ := utils.HashPassword(u.Password) // Hash the password
    u.Password = h
    config.DB.Create(&u) // Create user in the database
    c.JSON(http.StatusCreated, u) // Return created user
}

// GetAllUsers retrieves all users from the database.
func GetAllUsers(c *gin.Context) {
    var users []models.User
    config.DB.Find(&users) // Find all users
    c.JSON(http.StatusOK, users) // Return users
}

// GetUserByID retrieves a user by ID.
func GetUserByID(c *gin.Context) {
    id := c.Param("id")
    key := "user:" + id
    ctx := context.Background()
    if data, err := config.RedisClient.Get(ctx, key).Result(); err == nil {
        c.Data(http.StatusOK, "application/json", []byte(data)) // Return user data from Redis cache
        return
    }
    var u models.User
    if config.DB.First(&u, id).Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}) // Return error if user not found
        return
    }
    j := fmt.Sprintf(`{"id":%d,"full_name":"%s","username":"%s","status":"%s","role":"%s","image_url":"%s"}`,
        u.ID, u.FullName, u.Username, u.Status, u.Role, u.ImageURL)
    config.RedisClient.Set(ctx, key, j, 0) // Cache user data in Redis
    c.JSON(http.StatusOK, u) // Return user
}

// UpdateUser updates user information.
func UpdateUser(c *gin.Context) {
    id := c.Param("id")
    var u models.User
    if config.DB.First(&u, id).Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}) // Return error if user not found
        return
    }
    var in models.User
    if c.BindJSON(&in) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // Return error for invalid input
        return
    }
    u.FullName = in.FullName
    u.Status = in.Status
    u.Role = in.Role
    config.DB.Save(&u) // Save updated user
    c.JSON(http.StatusOK, u) // Return updated user
}

// DeleteUser deletes a user by ID.
func DeleteUser(c *gin.Context) {
    id := c.Param("id")
    config.DB.Delete(&models.User{}, id) // Delete user from database
    config.RedisClient.Del(context.Background(), "user:"+id) // Remove user from Redis cache
    c.JSON(http.StatusOK, gin.H{"deleted": id}) // Return deleted user ID
}

// UploadImage uploads a user image.
func UploadImage(c *gin.Context) {
    id := c.Param("id")
    var u models.User
    if config.DB.First(&u, id).Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}) // Return error if user not found
        return
    }
    f, h, err := c.Request.FormFile("image") // Get image from request
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No image"}) // Return error if no image provided
        return
    }
    defer f.Close()
    os.MkdirAll("uploads", os.ModePerm) // Create uploads directory
    name := fmt.Sprintf("uploads/%s_%s", u.Username, h.Filename) // Create image file name
    out, _ := os.Create(name) // Create image file
    defer out.Close()
    io.Copy(out, f) // Save image to file
    u.ImageURL = "/" + name
    config.DB.Save(&u) // Save updated user with image URL
    c.JSON(http.StatusOK, gin.H{"image_url": u.ImageURL}) // Return image URL
}

// DeleteImage deletes a user's image.
func DeleteImage(c *gin.Context) {
    id := c.Param("id")
    var u models.User
    if config.DB.First(&u, id).Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}) // Return error if user not found
        return
    }
    if u.ImageURL != "" {
        os.Remove(strings.TrimPrefix(u.ImageURL, "/")) // Remove image file
        u.ImageURL = ""
        config.DB.Save(&u) // Save updated user
    }
    c.JSON(http.StatusOK, gin.H{"deleted": id}) // Return deleted user ID
}

// isAlpha checks if a string contains only alphabetic characters.
func isAlpha(s string) bool {
    for _, c := range s {
        if c < 'A' || (c > 'Z' && c < 'a') || c > 'z' {
            return false // Return false if non-alphabetic character found
        }
    }
    return true // Return true if all characters are alphabetic
}
