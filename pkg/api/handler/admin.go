package handler

import (
	"fmt"
	"log"
	"net/http"
	"readon/pkg/api/middleware"
	domain "readon/pkg/domain"
	"readon/pkg/models"
	services "readon/pkg/usecase/interface"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type AdminHandler struct {
	AdminUseCase services.AdminUseCase
}

func NewAdminHandler(usecase services.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		AdminUseCase: usecase,
	}
}

// @Summary Get login page for admin
// @Description Get the HTML page for admin login
// @Tags Admin
// @Produce html
// @Success 200 {string} string "got html page: login as admin"
// @Router /adminlogin [get]
func (cr AdminHandler) GetLogin(c *gin.Context) {
	c.JSON(http.StatusOK, "got html page : login as admin")
}

// @Summary Admin Login
// @Description Login as an admin
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param admin body models.Userlogindata true "Admin credentials"
// @Success 200 {string} string "logged in TOKEN_STRING"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Router /adminlogin [post]
func (cr AdminHandler) Login(c *gin.Context) {
	var admin models.LoginData
	err := c.Bind(&admin)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    "Bad Request",
			Status: "error while binding form data",
			Hint:   "please check your request payload",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	id, is_admin := cr.AdminUseCase.Login(admin)

	if !is_admin {
		errResponse := models.ErrorResponse{
			Err:    "Unauthorised !",
			Status: "match not found , could not login :(",
			Hint:   "please try again",
		}
		c.JSON(http.StatusUnauthorized, errResponse)
		return
	}
	tokenString := middleware.GetTokenString(uint(id), "admin", false)

	response := "logged in  || TOKENSTRING : " + tokenString

	c.JSON(http.StatusOK, response)
	// setting cookie  here >>>
	// to do : should locate the setting cookie into  middleware
	c.SetCookie("Authorise", tokenString, 3600, "", "", true, false)
}

// @Summary List Admins
// @Description Get a list of Admins
// @Tags Admin
// @Produce json
// @Success 200 {array} models.Admin "List of Admins"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /admin/admins [get]
func (cr *AdminHandler) ListAdmins(c *gin.Context) {
	var list []models.Admin
	list, err := cr.AdminUseCase.ListAdmins()
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "InternalServerError",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, list)
}

// ListUsers godoc
// @summary Get all users
// @description Lists all users from the database.
// @description jwt temp:  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwicm9sZSI6ImFkbWluIn0.bw-k7fjfW6nY9DIBZ46ZcJG0IdaYclwHW7P7IPwuNoQ
// @tags Admin
// @produce json
// @Router /admin/users [post]
// @Param page query int false "Page number for pagination (default: 1)"
// @Param search query string false "Search word if any"
// @Security Bearer
// @Security Cookie
// @Success 200 {object} models.UselistResponse "OK"
// @Failure 400 {object} models.ErrorResponse "BadRequest"
// @Failure 500 {object} models.ErrorResponse "InternalServerError"
func (cr AdminHandler) ListUsers(c *gin.Context) {

	var pagedetails models.Pagination
	err := c.BindQuery(&pagedetails)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while parsing data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	log.Println("query got :", pagedetails)

	users, numOfResult, err := cr.AdminUseCase.ListUsers(pagedetails)
	log.Println("total num of results", numOfResult)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "couldn't get the users",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	log.Println("page detail", pagedetails)

	response := models.UserlistResponse{
		Pagination: pagedetails,
		List:       users,
	}

	c.JSON(http.StatusOK, response)

}

// @Summary Find user by ID
// @Description Get a user by ID
// @Tags Admin
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} models.User "User details"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Router /admin/user/{id} [get]
func (cr *AdminHandler) FindByID(c *gin.Context) {
	paramsId := c.Param("userId")
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

	user, err := cr.AdminUseCase.FindByID(uint(id))

	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "User doesn't exist !",
			Hint:   "please try again !",
		}
		c.JSON(http.StatusNotFound, errResponse)
		return
	} else {
		response := models.User{}
		copier.Copy(&response, &user)

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Delete user by ID
// @Description Delete a user by ID
// @Tags Admin
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {string} string "User is deleted successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /admin/user/{id} [delete]
func (cr *AdminHandler) Delete(c *gin.Context) {
	paramsId := c.Param("userId")
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

	//ctx := c.Request.Context()
	user, err := cr.AdminUseCase.FindByID(uint(id))

	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "couldn't find user ! ",
			Hint:   "please try again !",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	if user == (domain.User{}) {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "User has not been booked yet",
			Hint:   "please try again !",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	cr.AdminUseCase.Delete(user)
	response := "User is deleted successfully"
	c.JSON(http.StatusOK, response)
}

// @Summary Block or unblock user by ID
// @Description Block or unblock a user by ID
// @Tags Admin
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {string} string "User is blocked/unblocked successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Router /admin/blockuser/{id} [put]
func (cr *AdminHandler) BlockOrUnBlock(c *gin.Context) {
	paramId := c.Param("userId")
	fmt.Println("parmID :", paramId)
	id, err := strconv.Atoi(paramId)

	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "BadRequest , couldn't parse id ",
			Hint:   "please try again !",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	var response string

	current_status := cr.AdminUseCase.BlockOrUnBlock(id)
	if current_status == true {
		response = "user is unblocked !"
	} else {
		response = "user is blocked !"
	}
	c.JSON(http.StatusOK, response)
}
