package handler

import (
	"fmt"
	"net/http"

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

type Response struct {
	ID       uint   `copier:"must"`
	Name     string `copier:"must"`
	Email    string `copier:"must"`
	Password string `copier:"must"`
}

func NewUserHandler(usecase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}

func (cr UserHandler) GetLogin(c *gin.Context) {
	c.JSON(http.StatusOK, "got html page : user login")
}

// FindAll godoc
// @summary Get all users
// @description Get all users
// @tags users
// @security ApiKeyAuth
// @id FindAll
// @produce json
// @Router /user/users [get]
// @response 200 {object} []Response "OK"
func (cr *UserHandler) FindAll(c *gin.Context) {

	users, err := cr.userUseCase.FindAll(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		response := []Response{}
		copier.Copy(&response, &users)

		c.JSON(http.StatusOK, response)
	}
}

//	Signup godoc
//
// @summary Get all users
// @description Get all users
// @tags users
// @security ApiKeyAuth
// @id Save User
// @produce json
// @Router /signup [Post]
// @response 200 {object} []Response "OK"
func (cr *UserHandler) SaveUser(c *gin.Context) {
	var user domain.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := cr.userUseCase.Save(c.Request.Context(), user)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		response := Response{}
		copier.Copy(&response, &user)

		c.JSON(http.StatusOK, response)
	}

}

func (cr *UserHandler) GetSignup(c *gin.Context) {
	Response := "succesfully got the html page"
	c.JSON(http.StatusOK, Response)
}

func (cr UserHandler) UserLogin(c *gin.Context) {

	var userinput models.Userlogindata
	err := c.Bind(&userinput)
	if err != nil {
		fmt.Println("error while binding form data :", err)
	}

	// ckecking the db to match given data and gets the the user id from db in return
	userid, is_user, err := cr.userUseCase.UserLogin(c.Request.Context(), userinput)
	fmt.Println("eroor :", err)
	if !is_user {
		c.JSON(http.StatusNotFound, gin.H{
			"userid": userid,
			"error ": err,
			"status": " could not login :(",
			"hint":   "please try again",
		})
		return
	}

	tokenString := middleware.GetTokenString(userid, "user")
	c.SetCookie("Authorization", tokenString, 3600, "", "", true, false)
	c.JSON(http.StatusOK, gin.H{
		"userid": userid,
		"status": "logged in",
		"token":  tokenString,
	})
	// setting cookie  here >>>

	fmt.Println("token string :", tokenString)
}

func (cr *UserHandler) UserHome(c *gin.Context) {
	Response := "succesfully got the user home page"
	c.JSON(http.StatusOK, Response)
}

// Otp login

func (cr *UserHandler) GetOtpLogin(c *gin.Context) {
	c.JSON(http.StatusOK, "got login page !!! enter emil")
}

func (cr *UserHandler) VerifyAndSendOtp(c *gin.Context) {
	var userinput models.Userlogindata
	err := c.Bind(&userinput)
	if err != nil {
		c.JSON(http.StatusBadRequest, "couldnot bind json")
		return
	}
	fmt.Println("email :", userinput)
	err = cr.userUseCase.VerifyAndSendOtp(c, userinput.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Unauthorised !",
			"err":    err.Error(),
			"hint":   "please try again eith another email",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "got email and otp is sent",
		"email":   userinput.Email,
		"hint":    "enter otp ",
		"button":  "verify",
		"button ": "resend otp",
	})
}

func (cr UserHandler) VerifyOtp(c *gin.Context) {
	var otpinput domain.Otp
	err := c.Bind(&otpinput)
	if err != nil {
		c.JSON(http.StatusBadRequest, "couldnot bind json")
		return
	}
	err = cr.userUseCase.VerifyOtp(otpinput.Otp, otpinput.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status ": "Unauthorised otp",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status ": "otp verified",
		"hint ":   "redirected into home page",
	})
}
