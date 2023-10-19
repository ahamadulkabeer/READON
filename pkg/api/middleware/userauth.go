package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserAuthorizationMiddleware(c *gin.Context) {
	fmt.Println("Authorising....")

	tokenstring := c.Request.Header.Get("Authorization")
	if tokenstring == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": "token not found"})
		c.Abort()
		return
	}

	token, err := validateToken(tokenstring)
	if err != nil {
		c.JSON(http.StatusSeeOther, gin.H{
			"status": "cookie not found ",
			"error":  err,
		})
		c.Abort()
		return
	}

	role, ok := checkRole(token)
	if !ok || role != "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized token"})
		c.Abort()
		return
	}

	id, ok := getID(token)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access id"})
		c.Abort()
		return
	}
	c.Set("id", id)
	c.Next()
}
