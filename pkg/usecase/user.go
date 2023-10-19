package usecase

import (
	"context"
	"errors"
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

func (c *userUseCase) FindAll(ctx context.Context) ([]domain.User, error) {
	users, err := c.userRepo.FindAll(ctx)
	return users, err
}

func (c *userUseCase) Save(ctx context.Context, user domain.User) (domain.User, error) {

	fmt.Println("user:", user)

	err := middleware.ValidateUserData(&user)
	if err != nil {
		return user, err
	}
	user, err = c.userRepo.Save(ctx, user)
	return user, err
}

func (c *userUseCase) UserLogin(ctx context.Context, userinput models.Userlogindata) (int, bool, error) {

	user, err := c.userRepo.FindByEmail(ctx, userinput.Email)
	if err != nil {
		return 0, false, err
	}
	if user.Password != userinput.Password {
		return 0, false, errors.New("password does not match")
	}
	if user.Permission == false {
		return 0, false, errors.New("user have benn blocked by the admin")
	}
	return int(user.ID), true, err
}

func (c userUseCase) VerifyAndSendOtp(ctx context.Context, email string) error {
	err := c.userRepo.CheckForEmail(ctx, email)
	if err != nil {
		return errors.New("there is no user with this email")
	}
	otp, err := middleware.GenerateAndSendOpt(email)
	fmt.Println("otp   : ", otp)
	if err != nil {
		return errors.New("could not send otp")
	}
	c.userRepo.SaveOtp(otp, email)
	return nil
}

func (c userUseCase) VerifyOtp(otp, email string) error {
	err := c.userRepo.VerifyOtp(otp, email)
	return err
}
