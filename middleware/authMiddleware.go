package middleware

import (
	"net/http"

	"github.com/MoulieshN/Go-JWT-Project.git/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		claims, msg := helpers.ValidateToken(clientToken)
		if msg != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			c.Abort()
			return
		}
		c.Set("Email", claims.Email)
		c.Set("FirstName", claims.FirstName)
		c.Set("LastName", claims.LastName)
		c.Set("Uid", claims.Uid)
		c.Set("UserType", claims.UserType)
		c.Next()
	}
}
