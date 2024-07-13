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
		c.JSON(http.StatusUnauthorized, responses.RespondWithError(http.StatusUnauthorized, "Please provide valid authentication credentials.", "Authorization header is missing."))
		c.Abort()
		return
	}
	// stripping 'bearer' from token string
	tokenString = strings.Split(tokenString, " ")[1]
	// validating token string
	token, err := validateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.RespondWithError(http.StatusUnauthorized, "Please provide valid authentication credentials.", err.Error()))
		c.Abort()
		return
	}

	claims := GetClaims(*token)

	role, ok := checkRole(claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.RespondWithError(http.StatusUnauthorized, "role not found in cookie", "request not authorised for this action"))
		c.Abort()
		return
	}

	if role != "user" {
		c.JSON(http.StatusForbidden, responses.RespondWithError(http.StatusForbidden, "not authorised for this action", "UnAuthorised"))
		c.Abort()
		return
	}

	id, ok := getID(claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.RespondWithError(http.StatusUnauthorized, "userId not found in cookie", "request not authorised for this action"))
		c.Abort()
		return
	}
	if id.(float64) <= 0 {
		c.JSON(http.StatusUnauthorized, responses.RespondWithError(http.StatusUnauthorized, "unexpected userId found in Authorization header", "unAuthorised"))
		c.Abort()
		return
	}

	c.Set("userId", int(id.(float64)))
	c.Next()
}
