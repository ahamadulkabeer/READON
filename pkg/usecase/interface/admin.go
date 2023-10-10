package interfaces

import (
	"context"
	domain "readon/pkg/domain"
	"readon/pkg/models"
)

type AdminUseCase interface {
	Login(context.Context, models.Userlogindata) (int, bool)
	ListUser(context.Context) ([]models.ListOfUser, error)
	FindByID(ctx context.Context, id uint) (domain.Users, error)
	Delete(ctx context.Context, user domain.Users) error
}
