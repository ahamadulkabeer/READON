package helpers

import (
	"errors"
	"fmt"
	"net/mail"
	"readon/pkg/models"
	"regexp"
)

// GO - VALIDATOR

var AlphaNumericRegex = "^[a-zA-Z0-9]+$"

var AlphaNumericPlusRegex = "^[a-zA-Z0-9@#$%^&*()_+-=[\\]{}|;:'\",.<>?!/\\\\]+$"

func validateAlphanumeric(value string) bool {
	return regexp.MustCompile(AlphaNumericRegex).MatchString(value)
}

func validateAlphanumericPlus(value string) bool {
	return regexp.MustCompile(AlphaNumericPlusRegex).MatchString(value)
}

func ValidateName(name string) (bool, error) {

	if len(name) < 4 {
		return false, errors.New("user name should be atleasst 4 letters")
	}

	if len(name) < 4 {
		return false, errors.New("user name shouldn't be more than 12 letters")
	}

	ok := validateAlphanumeric(name)
	if !ok {
		return false, errors.New("user name should be alphanumeric")
	}
	return true, nil
}

func ValidatePassword(password string) (bool, error) {
	if len(password) < 8 || len(password) > 16 {
		return false, errors.New("password should should be between 8 to 16 letters")
	}
	ok := validateAlphanumericPlus(password)
	if !ok {
		return false, errors.New("user name should be alphanumericPlus")
	}
	return true, nil
}

func ValidateEmail(email string) (bool, error) {
	if validateAlphanumericPlus(email) {
		_, err := mail.ParseAddress(email)
		if err == nil {
			return true, nil
		}
	}
	return false, errors.New("invalid email")
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

func ValidateUserUPdateData(user *models.UserUpdateData) (bool, error) {

	return ValidateName(user.Name)

}
