package middleware

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func GetTokenString(id int, role string, premium bool) string {

	claims := jwt.MapClaims{
		"id":      id,
		"role":    role,
		"premium": premium,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println("error while creating token:", err)
	}
	return ss
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

func checkRole(token *jwt.Token) (string, bool) {
	claims, ok := token.Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	return role, ok
}

func getID(token *jwt.Token) (int, bool) {
	claims, ok := token.Claims.(jwt.MapClaims)
	id := int(claims["id"].(float64))
	return id, ok
}
