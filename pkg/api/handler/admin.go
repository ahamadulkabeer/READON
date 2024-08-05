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

func (cr AdminHandler) GetLogin(c *gin.Context) {
	c.JSON(http.StatusOK, "got html page : login as admin")
}

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

func (cr *AdminHandler) ListAdmins(c *gin.Context) {

	// pagination : ???
	response := cr.AdminUseCase.ListAdmins()

	c.JSON(response.StatusCode, response)
}

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
