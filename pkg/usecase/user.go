package usecase

import (
	"errors"
	"fmt"
	"net/http"

	"readon/pkg/api/helpers"
	"readon/pkg/api/responses"
	domain "readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"

	"github.com/jinzhu/copier"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(urepo interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: urepo,
	}
}

func (c *userUseCase) Save(userInput models.SignupData) responses.Response {

	fmt.Println("user:", userInput)

	var newUser domain.User
	copier.Copy(&newUser, &userInput)

	errors := struct {
		UserNameErr string
		EmailErr    string
		PasswordErr string
		GeneralErr  string
	}{}
	validated := true
	ok, err := helpers.ValidateName(newUser.Name)
	if !ok {
		validated = false
		errors.UserNameErr = err.Error()
	}
	ok, err = helpers.ValidateEmail(newUser.Email)
	if !ok {
		validated = false
		errors.EmailErr = err.Error()
	}
	ok, err = helpers.ValidatePassword(newUser.Password)
	if !ok {
		validated = false
		errors.PasswordErr = err.Error()
	}
	if !validated {
		return responses.ClientReponse(http.StatusBadRequest,
			"user validation failed", errors, nil)
	}

	exist, err := c.userRepo.CheckForEmail(userInput.Email)
	if err != nil {
		errors.GeneralErr = err.Error()
		return responses.ClientReponse(http.StatusInternalServerError,
			"Error while checking for email", errors, newUser)
	}
	if exist {
		errors.EmailErr = "email already in use "
		return responses.ClientReponse(http.StatusUnprocessableEntity,
			"invalid Email", errors, newUser)
	}
	userInput.Password, err = helpers.HashPassword(userInput.Password)
	if err != nil {
		errors.GeneralErr = "error while hashing password"
		return responses.ClientReponse(http.StatusInternalServerError,
			"couldn't process the request , please try again", errors, newUser)
	}

	newUser, err = c.userRepo.Save(newUser)
	if err != nil {
		errors.GeneralErr = err.Error()
		return responses.ClientReponse(http.StatusInternalServerError,
			"couldnt create user ", errors, newUser)
	}
	return responses.ClientReponse(http.StatusCreated,
		"user created succefully", nil, newUser)
}

func (c *userUseCase) UpdateUser(user models.UserUpdateData) (domain.User, error) {

	fmt.Println("user:", user)
	var User domain.User
	copier.Copy(&User, &user)

	err := helpers.ValidateUserUPdateData(&user)
	if err != nil {
		return User, err
	}
	User, err = c.userRepo.UpdateUser(User)
	return User, err
}

func (c *userUseCase) UserLogin(userinput models.LoginData) responses.Response {

	user, err := c.userRepo.FindByEmail(userinput.Email)
	if err != nil {
		return responses.ClientReponse(http.StatusInternalServerError, "couldn't login ", "error while searching for email", nil)
	}
	if user.ID == 0 {
		fmt.Println(" Email not found in db")
		return responses.ClientReponse(http.StatusNotFound, "couldn't login ", "invalid email or password", nil)
	}
	passed, _ := helpers.AuthenticatePassword(user.Password, userinput.Password)
	// if err != nil {
	// 	return responses.ClientReponse(http.StatusInternalServerError, "couldn't login ", "error while decrypting password", nil)
	// }
	if !passed {
		return responses.ClientReponse(http.StatusNotFound, "couldn't login ", "invalid email or password", nil)
	}
	if !user.Permission {
		fmt.Println("use have bee blocked !")
		return responses.ClientReponse(http.StatusForbidden, "couldn't login ", "user have been banned by the admin", nil)
	}

	return responses.ClientReponse(http.StatusOK, "user veryfied : redirecting ... ", nil, user)

}

func (c userUseCase) GetUserProfile(id int) (domain.User, error) {
	var user domain.User
	user, err := c.userRepo.FindByID(uint(id))
	if err != nil {
		return user, err
	}
	return user, err
}

func (c userUseCase) DeleteUserAccount(id int) error {
	var user domain.User
	user, err := c.userRepo.FindByID(uint(id))
	if err != nil {
		return err
	}
	err = c.userRepo.DeleteUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (c userUseCase) VerifyAndSendOtp(email string) error {
	_, err := c.userRepo.CheckForEmail(email)
	if err != nil {
		return errors.New("Invalid email")
	}
	otp, err := helpers.GenerateAndSendOpt(email)
	fmt.Println("otp   : ", otp)
	if err != nil {
		return errors.New("could not send otp")
	}
	c.userRepo.SaveOtp(otp, email)
	return nil
}

func (c userUseCase) VerifyOtp(otp, email string) error {
	err := c.userRepo.VerifyOtp(otp, email)
	return err
}
