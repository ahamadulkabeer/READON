package usecase

import (
	"context"
	"fmt"
	domain "readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
)

type AdminUseCase struct {
	adminRepo interfaces.AdminRepository
}

func NewAdminUsecase(repo interfaces.AdminRepository) services.AdminUseCase {
	return &AdminUseCase{
		adminRepo: repo,
	}
}

func (c AdminUseCase) Login(ctx context.Context, admin models.Userlogindata) (int, bool) {
	return c.adminRepo.Login(ctx, admin)
}

func (cr *AdminUseCase) ListUser(ctx context.Context) ([]models.ListOfUser, error) {
	return cr.adminRepo.ListUser(ctx)
}

func (c *AdminUseCase) FindByID(ctx context.Context, id uint) (domain.User, error) {
	user, err := c.adminRepo.FindByID(ctx, id)
	return user, err
}

func (c *AdminUseCase) Delete(ctx context.Context, user domain.User) error {
	err := c.adminRepo.Delete(ctx, user)
	return err
}

func (c *AdminUseCase) BlockOrUnBlock(ctx context.Context, id int) bool {
	status := c.adminRepo.BlockOrUnBlock(ctx, id)
	fmt.Println("status in usecase :", status)
	return status
}
