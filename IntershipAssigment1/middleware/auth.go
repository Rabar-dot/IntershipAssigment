package middleware

import (
    "fmt"
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware checks for a valid JWT token in the Authorization header.
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        h := c.GetHeader("Authorization") // Get Authorization header
        if h == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"}) // Return error if header is missing
            c.Abort()
            return
        }
        tokenString := strings.TrimPrefix(h, "Bearer ") // Extract token from header
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil // Return secret key for token validation
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"}) // Return error if token is invalid
            c.Abort()
            return
        }
        claims := token.Claims.(jwt.MapClaims) // Extract claims from token
        c.Set("username", claims["username"]) // Set username in context
        c.Set("role", claims["role"]) // Set role in context
        c.Next() // Proceed to the next handler
    }
}

// AuthorizeRole checks if the user has the required role.
func AuthorizeRole(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetString("role") != role {
            c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Requires %s role", role)}) // Return error if role does not match
            c.Abort()
            return
        }
        c.Next() // Proceed to the next handler
    }
}
