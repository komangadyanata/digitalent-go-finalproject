package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GetContentType(c *gin.Context) string {
	return c.Request.Header.Get("Content-Type")
}

func GetUserIdFromToken(c *gin.Context) uint {
	userData := c.MustGet("userData").(jwt.MapClaims)
	return uint(userData["id"].(float64))
}
