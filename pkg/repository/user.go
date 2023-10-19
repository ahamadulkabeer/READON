package repository

import (
	"context"
	"errors"
	"fmt"

	domain "readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"

	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

// Initilising repository

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB}
}

func (c *userDatabase) FindAll(ctx context.Context) ([]domain.User, error) {
	fmt.Println("ctx :", ctx)
	var users []domain.User
	err := c.DB.Find(&users).Error

	return users, err
}

func (c *userDatabase) Save(ctx context.Context, user domain.User) (domain.User, error) {
	err := c.DB.Save(&user).Error

	return user, err
}

func (c *userDatabase) Authorise(ctx context.Context, user models.Userlogindata) (int, bool, error) {
	var users domain.User
	result := c.DB.Where("email = ? AND password = ? AND permission = ?", user.Email, user.Password, true).Limit(1).Find(&users)
	if result.Error != nil {
		fmt.Println("error while ckecking for mathing data :", result.Error)
	}

	if result.RowsAffected == 0 {
		return 0, false, nil
	}
	return int(users.ID), true, result.Error
}

func (c *userDatabase) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := c.DB.Where("email = ?", email).Find(&user)
	if err != nil {
		return user, err.Error
	}
	return user, err.Error
}

func (c userDatabase) CheckForEmail(ctx context.Context, email string) error {
	var user domain.User
	result := c.DB.Where("email = ?", email).First(&user)
	return result.Error
}

func (c userDatabase) SaveOtp(otp, email string) error {
	var tosave = domain.Otp{Email: email, Otp: otp}
	result := c.DB.Save(&tosave)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c userDatabase) VerifyOtp(otp string, email string) error {
	var tocheck domain.Otp
	result := c.DB.Where("otp = ? AND email = ? ", otp, email).Limit(1).Find(&tocheck)
	if result.Error != nil {
		return errors.New("error from db !!!")
	}
	if result.RowsAffected == 0 {
		return errors.New("otp not found ")
	}
	return nil
}
