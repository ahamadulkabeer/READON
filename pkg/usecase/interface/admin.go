package interfaces

import (
	"context"
	domain "readon/pkg/domain"
	"readon/pkg/models"
)

type AdminUseCase interface {
	Login(context.Context, models.Userlogindata) (int, bool)
	ListUser(context.Context) ([]models.ListOfUser, error)
	FindByID(ctx context.Context, id uint) (domain.User, error)
	Delete(ctx context.Context, user domain.User) error
	BlockOrUnBlock(context.Context, int) bool
}
