package helpers

import (
	"errors"
	"fmt"
	"net/mail"
	"readon/pkg/domain"
	"readon/pkg/models"
	"regexp"

	validator "github.com/go-playground/validator/v10"
)

// GO - VALIDATOR
// should be moved the middleware  i think
var validate = validator.New()

func ValidateUserData(user *domain.User) error {

	validate.RegisterValidation("name", validateName)
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("email", validateEmail)
	if err := validate.Struct(user); err != nil {
		// Validation failed
		fmt.Println("Validation Error:", err)
		return err
	}

	// Validation passed
	fmt.Println("Validation successful")
	return nil
}

func ValidateUserUPdateData(user *models.UserUpdateData) error {

	validate.RegisterValidation("name", validateName)

	if err := validate.Struct(user); err != nil {
		// Validation failed
		fmt.Println("Validation Error:", err)
		return err
	}

	// Validation passed
	fmt.Println("Validation successful")
	return nil
}

func ValidateCategory(category string) error {

	if category == "" {
		return errors.New("category must not be empty")
	}
	if len(category) < 2 || len(category) > 20 {
		return errors.New("category must be atleast 2 letters ; max 20")
	}
	match, err := regexp.MatchString("^[a-zA-Z0-9]*$", category)
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("category must contain only letters and digits")
	}
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
		return validateAlphanumericPlus(password)
	}
	return false
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if validateAlphanumericPlus(email) {
		_, err := mail.ParseAddress(email)
		if err == nil {
			fmt.Println("here email return true ?????/")
			return true
		}
	}
	return false
}
func validateAlphanumeric(value string) bool {
	alphaNumericRegex := "^[a-zA-Z0-9]+$"
	return regexp.MustCompile(alphaNumericRegex).MatchString(value)
}

func validateAlphanumericPlus(value string) bool {
	alphaNumericRegex := "^[a-zA-Z0-9@#$%^&*()_+-=[\\]{}|;:'\",.<>?!/\\\\]+$"
	return regexp.MustCompile(alphaNumericRegex).MatchString(value)
}
