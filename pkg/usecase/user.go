package usecase

import (
	"context"
	"fmt"

	"readon/pkg/api/middleware"
	domain "readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (c *userUseCase) FindAll(ctx context.Context) ([]domain.Users, error) {
	users, err := c.userRepo.FindAll(ctx)
	return users, err
}

func (c *userUseCase) Save(ctx context.Context, user domain.Users) (domain.Users, error) {

	fmt.Println("user:", user)

	err := middleware.ValidateUserData(&user)
	if err != nil {
		return user, err
	}
	user, err = c.userRepo.Save(ctx, user)
	return user, err
}

func (c *userUseCase) UserLogin(ctx context.Context, userinput models.Userlogindata) (int, bool) {
	return c.userRepo.Authorise(ctx, userinput)
}
