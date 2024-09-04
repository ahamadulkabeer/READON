package usecase

import (
	"fmt"
	"net/http"

	"readon/pkg/api/errorhandler"
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
	// initialising and populates domain models user
	var newUser domain.User
	copier.Copy(&newUser, &userInput)

	// initialising custom error for errors
	errors := models.UserDataError{}

	// validating user input
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
			"user validation failed", errors, userInput)
	}

	// check if the email already in use
	exist, err := c.userRepo.CheckForEmail(userInput.Email)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		errors.GeneralErr = err.Error()
		return responses.ClientReponse(statusCode,
			"Error while checking for email", errors, userInput)
	}
	if exist {
		errors.EmailErr = "email already in use "
		return responses.ClientReponse(http.StatusUnprocessableEntity,
			"invalid Email", errors, userInput)
	}

	// hash password
	newUser.Password, err = helpers.HashPassword(userInput.Password)
	if err != nil {
		errors.GeneralErr = "error while hashing password"
		return responses.ClientReponse(http.StatusInternalServerError,
			"couldn't process the request , please try again", errors, userInput)
	}

	// create user
	newUser, err = c.userRepo.Save(newUser)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		errors.GeneralErr = err.Error()
		return responses.ClientReponse(statusCode,
			"couldnt create user ", errors, userInput)
	}

	// response with nessessary data
	var userData models.User
	copier.Copy(&userData, &newUser)
	return responses.ClientReponse(http.StatusCreated,
		"user created succefully", nil, userData)
}

func (c *userUseCase) UpdateUser(user models.UserUpdateData) responses.Response {
	//initialise the doamin models user
	var User domain.User
	copier.Copy(&User, &user)

	// validate user update data
	ok, err := helpers.ValidateUserUPdateData(&user)
	if err != nil || !ok {
		return responses.ClientReponse(http.StatusBadRequest, "couldn't update user data", "validation failed", nil)
	}

	// update user data
	User, err = c.userRepo.UpdateUser(User)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't update user data", err.Error(), err)
	}

	// response with nessessary data
	var userProfile models.User
	copier.Copy(&userProfile, &user)
	return responses.ClientReponse(http.StatusOK, "user details updated", nil, userProfile)

}

func (c *userUseCase) UserLogin(userinput models.LoginData) responses.Response {

	// check if the email exist
	user, err := c.userRepo.FindByEmail(userinput.Email)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode,
			"couldn't login ", "error while searching for email : "+err.Error(), userinput)
	}
	if user.ID == 0 {
		return responses.ClientReponse(http.StatusNotFound,
			"couldn't login ", "invalid email or password", userinput)
	}

	// if email and password match
	passed, err := helpers.AuthenticatePassword(user.Password, userinput.Password)
	if err != nil {
		return responses.ClientReponse(http.StatusUnprocessableEntity,
			"couldn't login ", err.Error(), userinput)
	}
	if !passed {
		return responses.ClientReponse(http.StatusNotFound, "couldn't login ",
			"invalid email or password", userinput)
	}

	// check if user have permission
	if !user.Permission {
		return responses.ClientReponse(http.StatusForbidden, "couldn't login ",
			"user have been blocked by the admin", userinput)
	}

	// response with nessessary data
	var userProfile models.User
	copier.Copy(&userProfile, &user)
	return responses.ClientReponse(http.StatusOK, "user logged in ", nil, userProfile)

}

func (c userUseCase) GetUserProfile(id int) responses.Response {
	// retrive user data
	user, err := c.userRepo.FindByID(uint(id))
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "coulndn't fetch user profile", err.Error(), nil)
	}

	// response
	var userProfile models.User
	copier.Copy(&userProfile, &user)
	return responses.ClientReponse(http.StatusOK, "user profile fetched ", nil, userProfile)
}

func (c userUseCase) DeleteUserAccount(id int) responses.Response {
	// check if user with the id exist
	user, err := c.userRepo.FindByID(uint(id))
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't delete user", err.Error(), nil)
	}
	// delete user
	err = c.userRepo.DeleteUser(user)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't delete user", err.Error(), nil)
	}
	// response
	var userProfile models.User
	copier.Copy(&userProfile, &user)
	return responses.ClientReponse(http.StatusOK, "user deleted successfully", nil, userProfile)
}

func (c userUseCase) VerifyAndSendOtp(email string) responses.Response {
	// exist, err := c.userRepo.CheckForEmail(email)
	// if err != nil {
	// 	statusCode, _ := errorhandler.HandleDatabaseError(err)
	// 	return responses.ClientReponse(statusCode,
	// 		"error finding email : otp not sent ", err.Error(), nil)
	// }
	// if !exist {
	// 	statusCode, _ := errorhandler.HandleDatabaseError(err)
	// 	return responses.ClientReponse(statusCode,
	// 		"no account found on email , please try again", "no account found", nil)
	// }
	otp, err := helpers.GenerateAndSendOpt(email)
	fmt.Println("otp   : ", otp)
	if err != nil {
		responses.ClientReponse(http.StatusOK, "couldn't generate otp ,please try again", err.Error(), nil)
	}
	err = c.userRepo.SaveOtp(otp, email)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode,
			"something went wrong , please try again", err.Error(), nil)
	}
	return responses.ClientReponse(http.StatusOK, "otp sent successfully", nil, "otp : "+otp)
}

func (c userUseCase) VerifyOtp(otpData domain.Otp) responses.Response {
	err := c.userRepo.VerifyOtp(otpData.Otp, otpData.Email)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		responses.ClientReponse(statusCode, "error verifying otp ,please try again", err.Error(), nil)
	}
	return responses.ClientReponse(http.StatusOK, "otp verified ", nil, nil)
}
