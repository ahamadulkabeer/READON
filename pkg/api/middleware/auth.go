package middleware

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

var JWTSecretKey string

func LoadJWTSecretKey(key string) error {
	if key == "" {
		return errors.New("jwt secret key is empty")
	}
	JWTSecretKey = key
	return nil
}

// creating token string from data
func GetTokenString(id uint, role string, premium bool) string {

	claims := jwt.MapClaims{
		"id":      id,
		"role":    role,
		"premium": premium,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(JWTSecretKey))
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

		return []byte(JWTSecretKey), nil
	})

	return token, err
}

func GetClaims(token jwt.Token) jwt.MapClaims {
	claims, _ := token.Claims.(jwt.MapClaims)
	return claims
}

func checkRole(claims jwt.MapClaims) (interface{}, bool) {
	role, ok := claims["role"]
	return role, ok
}

func getID(claims jwt.MapClaims) (interface{}, bool) {
	id, ok := claims["id"]
	return id, ok
}
