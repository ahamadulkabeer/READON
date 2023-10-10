package repository

import (
	"context"
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

func (c *userDatabase) FindAll(ctx context.Context) ([]domain.Users, error) {
	fmt.Println("ctx :", ctx)
	var users []domain.Users
	err := c.DB.Find(&users).Error

	return users, err
}

func (c *userDatabase) Save(ctx context.Context, user domain.Users) (domain.Users, error) {
	err := c.DB.Save(&user).Error

	return user, err
}

func (c *userDatabase) Authorise(ctx context.Context, user models.Userlogindata) (int, bool) {
	var users domain.Users
	result := c.DB.Where("email = ? AND password = ?", user.Email, user.Password).Limit(1).Find(&users)
	if result.Error != nil {
		fmt.Println("error while ckecking for mathing data :", result.Error)
	}

	if result.RowsAffected == 0 {
		return 0, false
	}
	return int(users.ID), true
}

// i dont know....:(
/*func (c *userDatabase) FindByEmail(ctx context.Context, email string) (domain.Users, error) {
	var user domain.Users
	err := c.DB.Where("email = ?", email).Find(&user)
	if err != nil {
		return user, err.Error
	}
	return user, err.Error
}*/
