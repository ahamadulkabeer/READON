package interfaces

import (
	"context"

	"readon/pkg/domain"
	"readon/pkg/models"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Save(ctx context.Context, user domain.User) (domain.User, error)
	Authorise(ctx context.Context, user models.Userlogindata) (int, bool, error)
	CheckForEmail(context.Context, string) error
	SaveOtp(string, string) error
	VerifyOtp(string, string) error
}
