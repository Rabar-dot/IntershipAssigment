package config

import (
	"fmt"
	"log"
	"os"

	"golang-user-api/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB               // Global variable for the database connection
var RedisClient *redis.Client // Global variable for the Redis client

// InitDB initializes the MySQL database connection.
func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")) // Data Source Name for MySQL

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // Open the database connection
	if err != nil {
		log.Fatal("Failed to connect MySQL:", err) // Log error if connection fails
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Migration failed:", err) // Log error if migration fails
	}
	DB = db
	fmt.Println("✅ MySQL connected and migrated.")
}

// InitRedis initializes the Redis client.
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // Redis address
		Password: os.Getenv("REDIS_PASS"), // Redis password
		DB:       0,                       // Default DB
	})
	if _, err := RedisClient.Ping(RedisClient.Context()).Result(); err != nil {
		log.Fatal("Failed to connect Redis:", err) // Log error if connection fails
	}
	fmt.Println("✅ Redis connected.")
}
