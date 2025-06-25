package middleware

import (
    "fmt"
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        h := c.GetHeader("Authorization")
        if h == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
            c.Abort()
            return
        }
        tokenString := strings.TrimPrefix(h, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        claims := token.Claims.(jwt.MapClaims)
        c.Set("username", claims["username"])
        c.Set("role", claims["role"])
        c.Next()
    }
}

func AuthorizeRole(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetString("role") != role {
            c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Requires %s role", role)})
            c.Abort()
            return
        }
        c.Next()
    }
}
