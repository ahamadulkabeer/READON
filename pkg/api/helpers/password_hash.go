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

func ValidatePassword(passwordStr, inputPassword string) (bool, error) {
	password := make([]byte, len(passwordStr))
	for i, x := range passwordStr {
		password[i] = byte(x)
	}
	fmt.Println(passwordStr == string(password))
	err := bcrypt.CompareHashAndPassword(password, []byte(inputPassword))
	if err != nil {
		fmt.Println("Password does not match")
		return false, err
	}
	return true, nil
}
