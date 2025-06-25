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

func CreateUser(c *gin.Context) {
    var u models.User
    if c.BindJSON(&u) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    if len(u.Username) == 0 || !isAlpha(u.Username) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
        return
    }
    if len(u.Password) < 8 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Password too short"})
        return
    }
    h, _ := utils.HashPassword(u.Password)
    u.Password = h
    config.DB.Create(&u)
    c.JSON(http.StatusCreated, u)
}

func GetAllUsers(c *gin.Context) {
    var users []models.User
    config.DB.Find(&users)
    c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
    id := c.Param("id")
    key := "user:" + id
    ctx := context.Background()
    if data, err := config.RedisClient.Get(ctx, key).Result(); err == nil {
        c.Data(http.StatusOK, "application/json", []byte(data))
        return
    }
    var u models.User
    if config.DB.First(&u, id).Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    j := fmt.Sprintf(`{"id":%d,"full_name":"%s","username":"%s","status":"%s","role":"%s","image_url":"%s"}`,
        u.ID, u.FullName, u.Username, u.Status, u.Role, u.ImageURL)
    config.RedisClient.Set(ctx, key, j, 0)
    c.JSON(http.StatusOK, u)
}

func UpdateUser(c *gin.Context) {
    id := c.Param("id")
    var u models.User
    if config.DB.First(&u, id).Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    var in models.User
    if c.BindJSON(&in) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    u.FullName = in.FullName
    u.Status = in.Status
    u.Role = in.Role
    config.DB.Save(&u)
    c.JSON(http.StatusOK, u)
}

func DeleteUser(c *gin.Context) {
    id := c.Param("id")
    config.DB.Delete(&models.User{}, id)
    config.RedisClient.Del(context.Background(), "user:"+id)
    c.JSON(http.StatusOK, gin.H{"deleted": id})
}

func UploadImage(c *gin.Context) {
    id := c.Param("id")
    var u models.User
    if config.DB.First(&u, id).Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    f, h, err := c.Request.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No image"})
        return
    }
    defer f.Close()
    os.MkdirAll("uploads", os.ModePerm)
    name := fmt.Sprintf("uploads/%s_%s", u.Username, h.Filename)
    out, _ := os.Create(name)
    defer out.Close()
    io.Copy(out, f)
    u.ImageURL = "/" + name
    config.DB.Save(&u)
    c.JSON(http.StatusOK, gin.H{"image_url": u.ImageURL})
}

func DeleteImage(c *gin.Context) {
    id := c.Param("id")
    var u models.User
    if config.DB.First(&u, id).Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    if u.ImageURL != "" {
        os.Remove(strings.TrimPrefix(u.ImageURL, "/"))
        u.ImageURL = ""
        config.DB.Save(&u)
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
