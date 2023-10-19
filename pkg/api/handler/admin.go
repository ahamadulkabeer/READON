package handler

import (
	"fmt"
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

func (cr AdminHandler) GetLogin(c *gin.Context) {
	c.JSON(http.StatusOK, "got html page : login as admin")
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

func (cr *AdminHandler) FindByID(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot parse id",
		})
		return
	}

	user, err := cr.AdminUseCase.FindByID(c.Request.Context(), uint(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User doesn't exist !",
		})
		return
	} else {
		response := Response{}
		copier.Copy(&response, &user)

		c.JSON(http.StatusOK, response)
	}
}

func (cr *AdminHandler) Delete(c *gin.Context) {
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
	user, err := cr.AdminUseCase.FindByID(ctx, uint(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User doesn't exist !",
		})
		return
		//c.AbortWithStatus(http.StatusNotFound)
	}

	if user == (domain.User{}) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User is not booking yet",
		})
		return
	}

	cr.AdminUseCase.Delete(ctx, user)

	c.JSON(http.StatusOK, gin.H{"message": "User is deleted successfully"})
}

func (cr *AdminHandler) BlockOrUnBlock(c *gin.Context) {
	paramId := c.Param("id")
	fmt.Println("parmID :", paramId)
	id, err := strconv.Atoi(paramId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot parse id",
		})
		return
	}
	current_status := cr.AdminUseCase.BlockOrUnBlock(c, id)
	if current_status == true {
		c.JSON(http.StatusOK, gin.H{"status": "user is unblocked !"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "user is blocked !"})
}
