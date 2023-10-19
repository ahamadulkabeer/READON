package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuthorizationMiddleware(c *gin.Context) {
	fmt.Println("Authorising....")

	tokenstring := c.Request.Header.Get("Authorization")
	if tokenstring == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": "tokenstring not found"})
		c.Abort()
		return
	}

	token, err := validateToken(tokenstring)

	if err != nil {
		c.JSON(http.StatusSeeOther, gin.H{
			"status": "invalid tokenstring",
			"error":  err,
		})
		c.Abort()
		return
	}

	role, ok := checkRole(token)
	if !ok || role != "admin" {
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

	ctx := context.WithValue(c, "userid", id)
	// Set the context to the request
	c.Request = c.Request.WithContext(ctx)
	c.Next()

}
