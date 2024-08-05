package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"readon/pkg/api/middleware"
	"readon/pkg/api/responses"
	models "readon/pkg/models"

	domain "readon/pkg/domain"
	services "readon/pkg/usecase/interface"
)

// initilising
type UserHandler struct {
	userUseCase services.UserUseCase
}

func NewUserHandler(usecase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}

// GetLogin godoc
// @Summary Get the login page
// @Description Retrieves the login page.
// @Tags Signup&Login
// @Produce json
// @Success 200 {object} responses.Response "login in page received"
// @Router /login [get]
func (cr UserHandler) GetLogin(c *gin.Context) {
	c.JSON(http.StatusOK, responses.ClientReponse(http.StatusOK, " login page recieved ", nil, nil))

}

// SaveUser godoc
// @Summary Save a new user
// @Description Save a new user to the system. Validates user data, checks for existing email, and hashes the password.
// @Tags Signup&Login
// @Accept json
// @Produce json
// @Param userInput formData models.SignupData true "User Input"
// @Success 201 {object} responses.Response{data=models.User} "Created"
// @Failure 400 {object} responses.Response{error=models.UserDataError,data=models.SignupData} "Bad Request"
// @Failure 422 {object} responses.Response{error=models.UserDataError,data=models.SignupData} "Unprocessable Entity"
// @Failure 500 {object} responses.Response{error=models.UserDataError,data=models.SignupData} "Internal Server Error"
// @Router /signup [post]
func (cr *UserHandler) SaveUser(c *gin.Context) {
	var user models.SignupData

	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error while bindin input data", err.Error(), nil))
		return
	}
	response := cr.userUseCase.Save(user)
	c.JSON(response.StatusCode, response)
}

// UpdateUser godoc
// @Summary Update user details
// @Description Update user details with the given data.user authentication required.
// @Tags User
// @Accept json
// @Produce json
// @Param user formData models.UserUpdateData true "User Update Data"
// @Success 200 {object} responses.Response{data=models.User} "User details updated"
// @Failure 400 {object} responses.Response "Bad Request"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /user/profile [put]
func (cr UserHandler) UpdateUser(c *gin.Context) {
	userId := c.GetInt("userId")
	var user models.UserUpdateData
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error while bindin data", err.Error(), nil))
		return
	}
	user.Id = userId
	response := cr.userUseCase.UpdateUser(user)

	c.JSON(response.StatusCode, response)
}

// GetSignup godoc
// @Summary Get the signup page
// @Description Retrieves the signup page.
// @Tags Signup&Login
// @Produce json
// @Success 200 {object} responses.Response "Sign up page received"
// @Router /signup [get]
func (cr *UserHandler) GetSignup(c *gin.Context) {
	c.JSON(http.StatusOK, responses.ClientReponse(http.StatusOK, " sign up page recieved ", nil, nil))

}

// UserLogin godoc
// @Summary User login
// @Description Log in a user with email and password.sets cookie valid for 24 hrs.
// @Tags Signup&Login
// @Accept json
// @Produce json
// @Param userInput formData models.LoginData true "User Login Data"
// @Success 200 {object} responses.Response{data=models.User} "User logged in"
// @Failure 400 {object} responses.Response "Bad Request"
// @Failure 404 {object} responses.Response "Not Found"
// @Failure 422 {object} responses.Response "Unprocessable Entity"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /login [post]
func (cr UserHandler) UserLogin(c *gin.Context) {

	var userInput models.LoginData

	if err := c.Bind(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error while bindin input data", err.Error(), nil))
		return
	}

	response := cr.userUseCase.UserLogin(userInput)
	if response.Error != nil {
		c.JSON(response.StatusCode, response)
		return
	}

	// settting cookie
	user := response.Data.(models.User)
	tokenString := middleware.GetTokenString(user.ID, "user", user.Permission)
	c.SetCookie("Authorization", tokenString, 3600, "", "", true, false)

	c.JSON(response.StatusCode, response)
}

// GetUserProfile godoc
// @Summary Get user profile
// @Description Fetch the user profile data for the logged-in user.User authentication required.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response{data=models.User} "User profile fetched"
// @Failure 404 {object} responses.Response{error=string} "Not Found"
// @Failure 500 {object} responses.Response{error=string} "Internal Server Error"
// @Router /user/profile [get]
func (cr UserHandler) GetUserProfile(c *gin.Context) {
	userId := c.GetInt("userId")
	response := cr.userUseCase.GetUserProfile(userId)

	c.JSON(response.StatusCode, response)
}

// DeleteUserAccount godoc
// @Summary Delete user account
// @Description Deletes the account of the logged-in user. User authentication required.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response{data=models.User} "User deleted successfully"
// @Failure 404 {object} responses.Response{error=string} "Not Found"
// @Failure 500 {object} responses.Response{error=string} "Internal Server Error"
// @Router /user/delete [delete]
func (cr UserHandler) DeleteUserAccount(c *gin.Context) {
	userId := c.GetInt("userId")
	response := cr.userUseCase.DeleteUserAccount(userId)
	c.JSON(response.StatusCode, response)

}

func (cr *UserHandler) UserHome(c *gin.Context) {
	Response := "succesfully got the user home page"
	c.JSON(http.StatusOK, Response)
}

// Otp login

func (cr *UserHandler) GetOtpLogin(c *gin.Context) {
	if c.GetHeader("Accept") == "application/json" {
		c.JSON(http.StatusOK, "got login page !!! enter emil")
		return
	}
	c.HTML(http.StatusOK, "login", nil)
}

func (cr *UserHandler) VerifyAndSendOtp(c *gin.Context) {

	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Could not process the request", "email is empty", nil))
		return
	}

	fmt.Println("email:", email)

	response := cr.userUseCase.VerifyAndSendOtp(email)

	c.JSON(response.StatusCode, response)
	//c.JSON(http.StatusOK, "Email verified, OTP is sent || Email :"+email)
}

func (cr *UserHandler) VerifyOtp(c *gin.Context) {
	var otpInput domain.Otp
	err := c.BindJSON(&otpInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Could not process the request", "binding error", nil))
		return
	}

	response := cr.userUseCase.VerifyOtp(otpInput)

	c.JSON(response.StatusCode, response)
}
