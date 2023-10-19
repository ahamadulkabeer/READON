package interfaces

import (
	"context"

	domain "readon/pkg/domain"
	"readon/pkg/models"
)

type UserUseCase interface {
	FindAll(ctx context.Context) ([]domain.User, error)

	Save(ctx context.Context, user domain.User) (domain.User, error)

	UserLogin(ctx context.Context, userinput models.Userlogindata) (int, bool, error)

	VerifyAndSendOtp(context.Context, string) error

	VerifyOtp(string, string) error
}
