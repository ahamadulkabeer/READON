package handler

import (
	"fmt"
	"net/http"
	"readon/pkg/api/middleware"
	"readon/pkg/models"
	services "readon/pkg/usecase/interface"

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

func (cr AdminHandler) Login(c *gin.Context) {
	var admin models.Userlogindata
	err := c.Bind(&admin)
	if err != nil {
		fmt.Println("error while binding form data :", err)
	}
	id, is_admin := cr.AdminUseCase.Login(c, admin)

	if !is_admin {
		c.JSON(http.StatusNotFound, gin.H{
			"userid": id,
			"status": "match not found , could not login :(",
			"hint":   "please try again",
		})
		return
	}
	tokenString := middleware.GetTokenString(id, "admin")

	c.JSON(http.StatusOK, gin.H{
		"userid": id,
		"status": "logged in",
		"token":  tokenString,
	})
	// setting cookie  here >>>

	fmt.Println("token string :", tokenString)

}

func (cr *AdminHandler) ListUser(c *gin.Context) {
	list, err := cr.AdminUseCase.ListUser(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "working",
		"users":  list,
	})
}
