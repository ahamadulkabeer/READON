package usecase

import (
	"net/http"
	"readon/pkg/api/errorhandler"
	"readon/pkg/api/middleware"
	"readon/pkg/api/responses"
	domain "readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"

	"github.com/gin-gonic/gin"
)

type AdminUseCase struct {
	adminRepo interfaces.AdminRepository
	userRepo  interfaces.UserRepository
}

func NewAdminUsecase(adminrepo interfaces.AdminRepository, userrepo interfaces.UserRepository) services.AdminUseCase {
	return &AdminUseCase{
		adminRepo: adminrepo,
		userRepo:  userrepo,
	}
}

func (c AdminUseCase) Login(admin models.LoginData, ctx *gin.Context) responses.Response {
	id, is_admin := c.adminRepo.Login(admin)

	if !is_admin {
		return responses.ClientReponse(http.StatusUnauthorized, "match not found , could not login :(", "Unauthorised !", nil)
	}

	tokenString := middleware.GetTokenString(uint(id), "admin", false)

	ctx.SetCookie("Authorise", tokenString, 3600, "", "", true, false)

	return responses.ClientReponse(http.StatusUnauthorized, "login successfull", nil, "TOKENSTRING : "+tokenString)
}

func (cr *AdminUseCase) ListAdmins() responses.Response {
	list, err := cr.adminRepo.ListAdmins()
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch list of admins", err, nil)
	}
	return responses.ClientReponse(http.StatusOK, "admin list fetched successfully", nil, list)
}

func (c AdminUseCase) ListUsers(pageDet models.Pagination) responses.Response {
	if pageDet.NewPage == 0 {
		pageDet.NewPage = 1
	}
	if pageDet.Size == 0 {
		pageDet.Size = 5
	}
	pageDet.Offset = pageDet.Size * (pageDet.NewPage - 1)
	users, numofresults, err := c.userRepo.ListUsers(pageDet)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch list of users data", err, nil)
	}
	// pageDet.Lastpage = numofresults / pageDet.Size
	// if numofresults%pageDet.Size != 0 {
	// 	pageDet.Lastpage++
	// }
	return responses.ClientReponse(http.StatusOK, "users data fetched successfully", nil, map[string]any{
		"currentpage":     pageDet.NewPage,
		"numberofresults": numofresults,
		"pagesize":        5,
		"data":            users,
	})
}

func (c *AdminUseCase) FindByID(id uint) responses.Response {
	user, err := c.userRepo.FindByID(id)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch user data", err, nil)
	}
	return responses.ClientReponse(http.StatusOK, "user data fetched successfully", nil, user)
}

func (c *AdminUseCase) Delete(user domain.User) responses.Response {
	err := c.userRepo.DeleteUser(user)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't delete user", err, nil)
	}
	return responses.ClientReponse(http.StatusOK, "user deleted successfully", nil, nil)
}

func (c *AdminUseCase) BlockOrUnBlock(id int) responses.Response {
	status, err := c.userRepo.BlockOrUnBlock(id)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't update user permission", err, nil)
	}

	return responses.ClientReponse(http.StatusOK, "user permission updated ", nil, map[string]any{
		"blocked": status,
	})
}
