package interfaces

import (
	"context"

	domain "readon/pkg/domain"
	"readon/pkg/models"
)

type UserUseCase interface {
	FindAll(ctx context.Context) ([]domain.Users, error)
	FindByID(ctx context.Context, id uint) (domain.Users, error)
	Save(ctx context.Context, user domain.Users) (domain.Users, error)
	Delete(ctx context.Context, user domain.Users) error
	UserLogin(ctx context.Context, userinput models.Userlogindata) (int, bool)
}
