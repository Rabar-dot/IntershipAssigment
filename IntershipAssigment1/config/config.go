package config

import (
    "fmt"
    "log"
    "os"

    "github.com/go-redis/redis/v8"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "golang-user-api/models"
)

var DB *gorm.DB
var RedisClient *redis.Client

func InitDB() {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"))

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect MySQL:", err)
    }
    if err := db.AutoMigrate(&models.User{}); err != nil {
        log.Fatal("Migration failed:", err)
    }
    DB = db
    fmt.Println("✅ MySQL connected and migrated.")
}

func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),
        Password: os.Getenv("REDIS_PASS"),
        DB:       0,
    })
    if _, err := RedisClient.Ping(RedisClient.Context()).Result(); err != nil {
        log.Fatal("Failed to connect Redis:", err)
    }
    fmt.Println("✅ Redis connected.")
}
