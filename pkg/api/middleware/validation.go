package middleware

import (
	"fmt"
	"readon/pkg/domain"
	"regexp"

	validator "github.com/go-playground/validator/v10"
)

// GO - VALIDATOR
// should be moved the middleware  i think
var validate = validator.New()

func ValidateUserData(user *domain.Users) error {

	validate.RegisterValidation("name", validateName)
	validate.RegisterValidation("password", validatePassword)

	if err := validate.Struct(user); err != nil {
		// Validation failed
		fmt.Println("Validation Error:", err)
		return err
	}

	// Validation passed
	fmt.Println("Validation successful")
	return nil

}

func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	// Check if the name is at least 4 characters long and  // ??contains no letters??
	if len(name) >= 4 {
		return validateAlphanumeric(name)
	}
	return false
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) >= 8 {
		return validateAlphanumeric(password)
	}
	return false
}

func validateAlphanumeric(value string) bool {
	//alphaNumericRegex := "^[a-zA-Z0-9]+$"
	alphaNumericRegex := "^[a-zA-Z0-9@#$%^&*()_+-=[\\]{}|;:'\",.<>?!/\\\\]+$"
	return regexp.MustCompile(alphaNumericRegex).MatchString(value)
}
