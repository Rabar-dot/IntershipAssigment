package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    FullName string `json:"full_name"`
    Username string `json:"username" gorm:"unique"`
    Password string `json:"-"`
    Status   string `json:"status"`
    Role     string `json:"role"`
    ImageURL string `json:"image_url"`
}
