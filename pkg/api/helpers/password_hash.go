package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

func AuthenticatePassword(password, inputPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(inputPassword))
	if err != nil {
		fmt.Println("Password does not match")
		return false, err
	}
	return true, nil
}
