package middleware

import (
	"fmt"
	"net/http"
	"readon/pkg/api/responses"
	"strings"

	"github.com/gin-gonic/gin"
)

func UserAuthorizationMiddleware(c *gin.Context) {
	fmt.Println("Authorising ... user")
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, responses.ClientReponse(http.StatusUnauthorized, "Please provide valid authentication credentials.", "Authorization header is missing.", nil))
		c.Abort()
		return
	}
	// stripping 'bearer' from token string
	tokenString = strings.Split(tokenString, " ")[1]
	// validating token string
	token, err := validateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.ClientReponse(http.StatusUnauthorized, "Please provide valid authentication credentials.", err.Error(), nil))
		c.Abort()
		return
	}

	claims := GetClaims(*token)

	role, ok := checkRole(claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.ClientReponse(http.StatusUnauthorized, "role not found in cookie", "request not authorised for this action", nil))
		c.Abort()
		return
	}

	if role != "user" {
		c.JSON(http.StatusForbidden, responses.ClientReponse(http.StatusForbidden, "not authorised for this action", "UnAuthorised", nil))
		c.Abort()
		return
	}

	id, ok := getID(claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.ClientReponse(http.StatusUnauthorized, "userId not found in cookie", "request not authorised for this action", nil))
		c.Abort()
		return
	}
	if id.(float64) <= 0 {
		c.JSON(http.StatusUnauthorized, responses.ClientReponse(http.StatusUnauthorized, "unexpected userId found in Authorization header", "unAuthorised", nil))
		c.Abort()
		return
	}

	c.Set("userId", int(id.(float64)))
	c.Next()
}
