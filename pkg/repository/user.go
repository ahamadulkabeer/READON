package repository

import (
	"errors"
	"fmt"
	"log"

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

func (c userDatabase) ListUsers(pageDet models.Pagination, offset int) ([]domain.User, int, error) {
	var users []domain.User
	var numOfResult int64
	log.Println("pagedat in repo", pageDet)
	log.Println("offset :", offset)
	err := c.DB.Table("users").Select("id,name,email,permission").Where(" name ILIKE  ?", fmt.Sprintf("%%%s%%", pageDet.Search)).Offset(offset).Limit(pageDet.Size).Find(&users).Error
	if err != nil {
		return users, 0, err
	}
	log.Println("users", users)
	err = c.DB.Table("users").Select("COUNT(*)").Where("name ILIKE ?", fmt.Sprintf("%%%s%%", pageDet.Search)).Count(&numOfResult).Error
	return users, int(numOfResult), err
}

func (c *userDatabase) Save(user domain.User) (domain.User, error) {
	err := c.DB.Save(&user).Error

	return user, err
}

func (c *userDatabase) UpdateUser(user domain.User) (domain.User, error) {
	err := c.DB.Model(&domain.User{}).Where("id = ?", user.ID).Update("name", user.Name).Error

	return user, err
}

func (c *userDatabase) Authorise(user models.Userlogindata) (int, bool, error) {
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

func (c userDatabase) FindByID(id uint) (domain.User, error) {
	var user domain.User
	err := c.DB.First(&user, id).Error

	return user, err
}

func (c userDatabase) DeleteUser(user domain.User) error {
	err := c.DB.Delete(&user).Error

	return err
}

func (c userDatabase) BlockOrUnBlock(id int) bool {
	sql := `
	        UPDATE users
	        SET permission = CASE
	            WHEN permission THEN false
	            ELSE true
	        END
	        WHERE id = ?`

	c.DB.Exec(sql, id)

	var permission bool
	sql = "SELECT permission FROM users WHERE id = ?"
	c.DB.Raw(sql, id).Scan(&permission)
	return permission
}

func (c *userDatabase) FindByEmail(email string) (domain.User, error) {
	var user domain.User
	db := c.DB.Where("email = ?", email).Find(&user)
	if db.RowsAffected <= 0 {
		return user, errors.New("Email not found")
	}
	log.Println("login user details : ", user)
	return user, nil
}

func (c userDatabase) CheckForEmail(email string) error {
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
