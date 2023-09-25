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

func (cr *UserHandler) FindByID(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot parse id",
		})
		return
	}

	user, err := cr.userUseCase.FindByID(c.Request.Context(), uint(id))

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		response := Response{}
		copier.Copy(&response, &user)

		c.JSON(http.StatusOK, response)
	}
}

func (cr *UserHandler) SaveUser(c *gin.Context) {
	var user domain.Users

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

func (cr *UserHandler) Delete(c *gin.Context) {
	paramsId := c.Param("id")
	fmt.Println("parmsID :", paramsId)
	id, err := strconv.Atoi(paramsId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot parse id",
		})
		return
	}

	ctx := c.Request.Context()
	user, err := cr.userUseCase.FindByID(ctx, uint(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User doesn't exist !",
		})
		return
		//c.AbortWithStatus(http.StatusNotFound)
	}

	if user == (domain.Users{}) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User is not booking yet",
		})
		return
	}

	cr.userUseCase.Delete(ctx, user)

	c.JSON(http.StatusOK, gin.H{"message": "User is deleted successfully"})
}

func (cr *UserHandler) GetSignup(c *gin.Context) {
	Response := "succesfully got the html page"
	c.JSON(http.StatusOK, Response)
}

func (cr UserHandler) UserLogin(c *gin.Context) {

	//email := c.PostForm("email")
	//password := c.PostForm("password")
	// instead
	var userinput models.Userlogindata
	err := c.Bind(&userinput)
	if err != nil {
		fmt.Println("error while binding form data :", err)
	}

	fmt.Println("email and password ", userinput)

	// ckecking the db to match given data and gets the the user id from db in return
	userid, is_user := cr.userUseCase.UserLogin(c.Request.Context(), userinput)
	if !is_user {
		c.JSON(http.StatusNotFound, gin.H{
			"userid": userid,
			"status": "match not found , could not login :(",
			"hint":   "please try again",
		})
		return
	}

	tokenString := middleware.GetTokenString(userid, "user")

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
