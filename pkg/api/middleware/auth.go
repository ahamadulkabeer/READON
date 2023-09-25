package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthorizationMiddleware(c *gin.Context) {
	fmt.Println("Authorising....")
	s := c.Request.Header.Get("Authorization")

	tokenstring := strings.TrimPrefix(s, "Bearer ")

	token, err := validateToken(tokenstring)
	if err != nil {
		c.Abort()
		fmt.Println("error in validating tokenstring :", err)
		c.JSON(http.StatusSeeOther, gin.H{
			"status": "cookie not found ,redirected to guest home",
			"error":  err,
		})
		return
	}
	role := ckeckingRole(token)

	switch role {
	case "admin":
		c.JSON(http.StatusOK, "got admin page")
	case "creator":
		c.JSON(http.StatusOK, "got creator home")
	case "user":
		c.JSON(http.StatusOK, "got user home")
	default:
		c.JSON(http.StatusNotFound, "couldn,t login ,reddirected to goes thomme")
	}

}

func ckeckingRole(token *jwt.Token) string {

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Access your claim data here
		role := claims["role"].(string)
		return role
		// You can access other claims in a similar wa
	}
	return ""
}

func validateToken(tokenstring string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte("secret"), nil
	})

	return token, err
}

func GetTokenString(id int, role string) string {

	claims := jwt.MapClaims{
		"id":   id,
		"role": role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println("error while creating token:", err)
	}
	return ss
}
