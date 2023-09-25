package usecase

import (
	"context"
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
