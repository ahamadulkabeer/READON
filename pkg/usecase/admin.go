package usecase

import (
	"fmt"
	domain "readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
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

func (c AdminUseCase) Login(admin models.Userlogindata) (int, bool) {
	return c.adminRepo.Login(admin)
}

func (cr *AdminUseCase) ListAdmins() ([]models.Admin, error) {
	return cr.adminRepo.ListAdmins()
}

func (c AdminUseCase) ListUsers(pageDet *models.Pagination) ([]domain.User, int, error) {
	if pageDet.NewPage == 0 {
		pageDet.NewPage = 1
	}
	pageDet.Size = 5
	offset := pageDet.Size * (pageDet.NewPage - 1)
	users, numofresults, err := c.userRepo.ListUsers(*pageDet, offset)
	pageDet.Lastpage = numofresults / pageDet.Size
	if numofresults%pageDet.Size != 0 {
		pageDet.Lastpage++
	}
	return users, numofresults, err
}

func (c *AdminUseCase) FindByID(id uint) (domain.User, error) {
	user, err := c.userRepo.FindByID(id)
	return user, err
}

func (c *AdminUseCase) Delete(user domain.User) error {
	err := c.userRepo.DeleteUser(user)
	return err
}

func (c *AdminUseCase) BlockOrUnBlock(id int) bool {
	status := c.userRepo.BlockOrUnBlock(id)
	fmt.Println("status in usecase :", status)
	return status
}
