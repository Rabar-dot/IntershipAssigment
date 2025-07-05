package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang-user-api/config" // Import config to use Logger
	"golang-user-api/models"
	"golang-user-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func CreateUser(c *gin.Context) {
	var u models.User
	if c.BindJSON(&u) != nil {
		config.Logger.WithField("error", "Invalid input").Error("Failed to bind JSON for CreateUser")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if len(u.Username) == 0 || !isAlpha(u.Username) {
		config.Logger.WithField("username", u.Username).Error("Invalid username for CreateUser")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
		return
	}
	if len(u.Password) < 8 {
		config.Logger.Error("Password too short for CreateUser")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password too short"})
		return
	}
	h, err := utils.HashPassword(u.Password)
	if err != nil {
		config.Logger.WithField("error", err.Error()).Error("Failed to hash password for CreateUser")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}
	u.Password = h
	if err := config.DB.Create(&u).Error; err != nil {
		config.Logger.WithField("error", err.Error()).Error("Failed to create user in DB")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	config.Logger.WithField("user_id", u.ID).Info("User created successfully")
	c.JSON(http.StatusCreated, u)
}

func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		config.Logger.WithField("error", err.Error()).Error("Failed to retrieve all users from DB")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	config.Logger.Info("Successfully retrieved all users")
	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	key := "user:" + id
	ctx := context.Background()

	// Try to get from Redis
	if data, err := config.RedisClient.Get(ctx, key).Result(); err == nil {
		config.Logger.WithField("user_id", id).Info("User found in Redis cache")
		c.Data(http.StatusOK, "application/json", []byte(data))
		return
	} else if err != redis.Nil && err != context.Canceled && err != context.DeadlineExceeded {
		config.Logger.WithField("user_id", id).WithField("error", err.Error()).Warn("Error getting user from Redis, trying DB")
	}

	var u models.User
	// Get from Database
	if err := config.DB.First(&u, id).Error; err != nil {
		config.Logger.WithField("user_id", id).WithField("error", err.Error()).Error("User not found in database")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Store in Redis
	j := fmt.Sprintf(`{"id":%d,"full_name":"%s","username":"%s","status":"%s","role":"%s","image_url":"%s"}`,
		u.ID, u.FullName, u.Username, u.Status, u.Role, u.ImageURL)
	if err := config.RedisClient.Set(ctx, key, j, 0).Err(); err != nil {
		config.Logger.WithField("user_id", id).WithField("error", err.Error()).Warn("Failed to set user in Redis cache")
	} else {
		config.Logger.WithField("user_id", id).Info("User stored in Redis cache")
	}

	config.Logger.WithField("user_id", id).Info("User retrieved from database")
	c.JSON(http.StatusOK, u)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var u models.User
	if err := config.DB.First(&u, id).Error; err != nil {
		config.Logger.WithField("user_id", id).WithField("error", err.Error()).Error("User not found for update")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var in models.User
	if c.BindJSON(&in) != nil {
		config.Logger.WithField("user_id", id).Error("Invalid input for UpdateUser")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	u.FullName = in.FullName
	u.Status = in.Status
	u.Role = in.Role
	if err := config.DB.Save(&u).Error; err != nil {
		config.Logger.WithField("user_id", id).WithField("error", err.Error()).Error("Failed to save user update to DB")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	// Invalidate Redis cache for the updated user
	config.RedisClient.Del(context.Background(), "user:"+id)
	config.Logger.WithField("user_id", id).Info("User updated and Redis cache invalidated")
	c.JSON(http.StatusOK, u)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		config.Logger.WithField("user_id", id).WithField("error", err.Error()).Error("Failed to delete user from DB")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	config.RedisClient.Del(context.Background(), "user:"+id)
	config.Logger.WithField("user_id", id).Info("User deleted and Redis cache invalidated")
	c.JSON(http.StatusOK, gin.H{"deleted": id})
}

func UploadImage(c *gin.Context) {
	id := c.Param("id")
	var u models.User
	if err := config.DB.First(&u, id).Error; err != nil {
		config.Logger.WithField("user_id", id).WithField("error", err.Error()).Error("User not found for image upload")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	f, h, err := c.Request.FormFile("image")
	if err != nil {
		config.Logger.WithField("user_id", id).WithField("error", err.Error()).Error("No image file provided for upload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image"})
		return
	}
	defer f.Close()
	os.MkdirAll("uploads", os.ModePerm)
	name := fmt.Sprintf("uploads/%s_%s", u.Username, h.Filename)
	out, err := os.Create(name)
	if err != nil {
		config.Logger.WithField("user_id", id).WithField("filename", name).WithField("error", err.Error()).Error("Failed to create image file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, f); err != nil {
		config.Logger.WithField("user_id", id).WithField("filename", name).WithField("error", err.Error()).Error("Failed to copy image data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}
	u.ImageURL = "/" + name
	if err := config.DB.Save(&u).Error; err != nil {
		config.Logger.WithField("user_id", id).WithField("error", err.Error()).Error("Failed to update user with image URL")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user image URL"})
		return
	}
	// Invalidate Redis cache for the updated user
	config.RedisClient.Del(context.Background(), "user:"+id)
	config.Logger.WithField("user_id", id).Info("Image uploaded and user updated, Redis cache invalidated")
	c.JSON(http.StatusOK, gin.H{"image_url": u.ImageURL})
}

func DeleteImage(c *gin.Context) {
	id := c.Param("id")
	var u models.User
	if err := config.DB.First(&u, id).Error; err != nil {
		config.Logger.WithField("user_id", id).WithField("error", err.Error()).Error("User not found for image deletion")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if u.ImageURL != "" {
		filePath := strings.TrimPrefix(u.ImageURL, "/")
		if err := os.Remove(filePath); err != nil {
			config.Logger.WithField("user_id", id).WithField("filepath", filePath).WithField("error", err.Error()).Warn("Failed to delete image file from disk")
			// Continue even if file deletion fails, as we still want to clear the URL in DB
		}
		u.ImageURL = ""
		if err := config.DB.Save(&u).Error; err != nil {
			config.Logger.WithField("user_id", id).WithField("error", err.Error()).Error("Failed to clear image URL in DB")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear image URL"})
			return
		}
		// Invalidate Redis cache for the updated user
		config.RedisClient.Del(context.Background(), "user:"+id)
		config.Logger.WithField("user_id", id).Info("Image URL cleared and Redis cache invalidated")
	} else {
		config.Logger.WithField("user_id", id).Info("No image URL to delete for user")
	}
	c.JSON(http.StatusOK, gin.H{"deleted": id})
}

func isAlpha(s string) bool {
	for _, c := range s {
		if c < 'A' || (c > 'Z' && c < 'a') || c > 'z' {
			return false
		}
	}
	return true
}
