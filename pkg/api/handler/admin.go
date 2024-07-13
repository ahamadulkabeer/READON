package handler

import (
	"log"
	"net/http"
	"readon/pkg/api/responses"
	"readon/pkg/models"
	services "readon/pkg/usecase/interface"
	"strconv"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest, "error while binding form data", err, nil))
		return
	}
	response := cr.AdminUseCase.Login(admin, c)

	c.JSON(response.StatusCode, response)
}

// @Summary List Admins
// @Description Get a list of Admins
// @Tags Admin
// @Produce json
// @Success 200 {array} models.Admin "List of Admins"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /admin/admins [get]
func (cr *AdminHandler) ListAdmins(c *gin.Context) {

	// pagination : ???
	response := cr.AdminUseCase.ListAdmins()

	c.JSON(response.StatusCode, response)
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

		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest, "Error while parsing data", err, nil))
		return
	}

	log.Println("query got :", pagedetails)

	response := cr.AdminUseCase.ListUsers(pagedetails)

	c.JSON(response.StatusCode, response)

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
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest, "Error while parsing data", err, nil))
		return
	}

	response := cr.AdminUseCase.FindByID(uint(id))

	c.JSON(response.StatusCode, response)
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
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest, "Error while parsing data", err, nil))
		return
	}

	response := cr.AdminUseCase.FindByID(uint(id))

	c.JSON(response.StatusCode, response)
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
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest, "Error while parsing data", err, nil))
		return
	}

	response := cr.AdminUseCase.BlockOrUnBlock(id)

	c.JSON(response.StatusCode, response)
}
