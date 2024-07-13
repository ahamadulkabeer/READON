package interfaces

import (
	"readon/pkg/api/responses"
	domain "readon/pkg/domain"
	"readon/pkg/models"

	"github.com/gin-gonic/gin"
)

type AdminUseCase interface {
	Login(admin models.LoginData, context *gin.Context) responses.Response
	ListAdmins() responses.Response
	ListUsers(models.Pagination) responses.Response
	FindByID(id uint) responses.Response
	Delete(user domain.User) responses.Response
	BlockOrUnBlock(int) responses.Response
}
