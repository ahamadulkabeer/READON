package interfaces

import (
	"context"

	"readon/pkg/domain"
	"readon/pkg/models"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]domain.Users, error)
	FindByID(ctx context.Context, id uint) (domain.Users, error)
	//FindByEmail(ctx context.Context, email string)(domain.Users, error)
	Save(ctx context.Context, user domain.Users) (domain.Users, error)
	Delete(ctx context.Context, user domain.Users) error
	Authorise(ctx context.Context, user models.Userlogindata) (int, bool)
}
