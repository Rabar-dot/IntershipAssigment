package utils

import (
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a new JWT token with the given username and role.
func GenerateJWT(username, role string) (string, error) {
    claims := jwt.MapClaims{
        "username": username,
        "role":     role,
        "exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Create a new token
    return token.SignedString([]byte(os.Getenv("JWT_SECRET"))) // Sign the token with the secret
}
