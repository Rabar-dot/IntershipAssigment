package models

import "gorm.io/gorm"

// User represents the user model in the database
// Make sure fields match the expected input/output in JSON and GORM

type User struct {
    gorm.Model
    FullName string `json:"full_name"`
    Username string `json:"username" gorm:"unique"`
    Password string `json:"password"`
    Status   string `json:"status"`
    Role     string `json:"role"`
    ImageURL string `json:"image_url"`
}