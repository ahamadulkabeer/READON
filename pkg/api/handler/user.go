package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"readon/pkg/api/middleware"
	models "readon/pkg/models"

	domain "readon/pkg/domain"
	services "readon/pkg/usecase/interface"
)

type UserHandler struct {
	userUseCase services.UserUseCase
}

// fo rsome reason could not get the struct from doamin domain.Otp had create one
type Otp struct {
	Email string `form:"email" json:"email"`
	Otp   string `form:"otp" json:"otp"`
}

func NewUserHandler(usecase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}

// GetLogin godoc
// @Summary Get HTML page for user login
// @Description Retrieve the HTML page for user login.
// @Tags users
// @Produce json
// @Success 200 {string} string "Got HTML page for user login"
// @Router /login [get]
func (cr UserHandler) GetLogin(c *gin.Context) {

	c.JSON(http.StatusOK, "got html page : user login")
}

// @Summary Save a user
// @Description Save a user by providing JSON payload
// @Tags users
// @Security ApiKeyAuth
// @ID SaveUser
// @Accept json
// @Produce json
// @Param usersdata body models.SignupData true "User object to be saved"
// @Router /signup [post]
// @Success 200 {object} domain.User "OK"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "InternalServerError"
func (cr *UserHandler) SaveUser(c *gin.Context) {
	var user models.SignupData

	if err := c.BindJSON(&user); err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "error while bindin json save user",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	User, err := cr.userUseCase.Save(user)

	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "couldn't Save users",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	} else {
		response := models.User{}
		copier.Copy(&response, &User)

		c.JSON(http.StatusOK, response)
	}
}

// GetSignup godoc
// @Summary Get HTML page for user signup
// @Description Retrieve the HTML page for user signup.
// @Tags users
// @Produce html
// @Success 200 {string} string "Successfully got the HTML page"
// @Router /signup [get]
func (cr *UserHandler) GetSignup(c *gin.Context) {
	Response := "succesfully got the html page"
	c.JSON(http.StatusOK, Response)
}

// UserLogin godoc
// @Summary User login
// @Description Log in a user and return a token if successful.
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.Userlogindata true "User login data"
// @Success 200 {string} string "succesfully logged in  + tokenstring "
// @Failure 401 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /login [post]
func (cr UserHandler) UserLogin(c *gin.Context) {

	var userinput models.Userlogindata
	err := c.Bind(&userinput)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Err:    "Bad Request",
			Status: "Could not process the request",
			Hint:   "Please check your request data",
		})
		return
	}
	// ckecking the db to match given data and gets the the user id from db in return
	userid, is_user, isPremium, err := cr.userUseCase.UserLogin(userinput)
	fmt.Println("eroor :", err)
	if !is_user {
		errresponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: " could not login :(",
			Hint:   "please try again",
		}
		c.JSON(http.StatusUnauthorized, errresponse)
		return
	}

	tokenString := middleware.GetTokenString(userid, "user", isPremium)
	c.SetCookie("Authorization", tokenString, 3600, "", "", true, false)
	Response := "successfully logged in :) || token : " + tokenString
	c.JSON(http.StatusOK, Response)

	fmt.Println("token string :", tokenString)
}

// @Summary Get user profile
// @Description Get user profile information by providing the user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID" Format(int64)
// @Success 200 {object} domain.User
// @Failure 400 {object} models.ErrorResponse "Bad Request, couldn't parse id"
// @Failure 500 {object} models.ErrorResponse "Couldn't get user profile"
// @Router /user/profile/{id} [get]
func (cr UserHandler) GetUserProfile(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)

	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "BadRequest , couldn't parse id ",
			Hint:   "please try again !",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	user, err := cr.userUseCase.GetUserProfile(id)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "couldn't delete user ! ",
			Hint:   "please try again !",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	response := user
	c.JSON(http.StatusOK, response)

}

// @Summary Delete a user account
// @Description Delete a user account by providing the user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID" Format(int64)
// @Success 200 {string} string "User is deleted successfully || Redirected to login page"
// @Failure 400 {object} models.ErrorResponse "Bad Request, couldn't parse id"
// @Failure 500 {object} models.ErrorResponse "Couldn't delete user"
// @Router /user/account/{id} [delete]
func (cr UserHandler) DeleteUserAccount(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)
	id = c.GetInt("id")
	id = c.MustGet("id").(int)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "BadRequest , couldn't parse id ",
			Hint:   "please try again !",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	err = cr.userUseCase.DeleteUserAccount(id)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "couldn't delete user ! ",
			Hint:   "please try again !",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	response := "User is deleted successfully  || " + "Redirected to login page"
	c.JSON(http.StatusOK, response)

}

// UserHome godoc
// @Summary Get user home page
// @Description Retrieve the user home page.
// @Tags users
// @Produce json
// @Success 200 {string} string "Successfully got the user home page"
// @Router /user/home [get]
func (cr *UserHandler) UserHome(c *gin.Context) {
	Response := "succesfully got the user home page"
	c.JSON(http.StatusOK, Response)
}

// Otp login

// GetOtpLogin godoc
// @Summary Get HTML page for OTP login
// @Description Retrieve the HTML page for OTP login.
// @Tags users
// @Produce json
// @Success 200 {string} string "Got HTML page for OTP login, enter email"
// @Router /otplogin [get]
func (cr *UserHandler) GetOtpLogin(c *gin.Context) {
	c.JSON(http.StatusOK, "got login page !!! enter emil")
}

// VerifyAndSendOtp godoc
// @Summary Verify email and send OTP
// @Description Verify the provided email and send an OTP.
// @Tags users
// @Accept json
// @Produce json
// @Param email body string true "User email to verify and send OTP"
// @Success 200 {string} string "Email verified, OTP sent  + Email"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Router /otplogin [post]
func (cr *UserHandler) VerifyAndSendOtp(c *gin.Context) {
	var email string

	if err := c.BindJSON(&email); err != nil {
		fmt.Println("email :", email)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Err:    "Bad Request",
			Status: "Could not process the request",
			Hint:   "Please check your request data",
		})
		return
	}

	fmt.Println("email:", email)

	err := cr.userUseCase.VerifyAndSendOtp(email)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Unauthorized!",
			Hint:   "Please try again with another email",
		}
		c.JSON(http.StatusUnauthorized, errResponse)
		return
	}

	c.JSON(http.StatusOK, "Email verified, OTP is sent || Email :"+email)
}

// VerifyOtp godoc
// @Summary Verify OTP
// @Description Verify the provided OTP.
// @Tags users
// @Accept json
// @Produce json
// @Param input body Otp true "OTP data"
// @Success 200 {string} string "OTP verified, redirected to home page"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized OTP"
// @Router /verifyotp [post]
func (cr *UserHandler) VerifyOtp(c *gin.Context) {
	var otpinput domain.Otp
	if err := c.BindJSON(&otpinput); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Err:    "Bad Request",
			Status: "Could not process the request",
			Hint:   "Please check your request data",
		})
		return
	}

	err := cr.userUseCase.VerifyOtp(otpinput.Otp, otpinput.Email)
	if err != nil {
		errresponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Invalid OTP!",
			Hint:   "Please try again",
		}
		c.JSON(http.StatusUnauthorized, errresponse)
		return
	}

	c.JSON(http.StatusOK, "OTP verified, redirected to the home page")
}
