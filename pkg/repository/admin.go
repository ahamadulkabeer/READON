package repository

import (
	"context"
	"fmt"
	domain "readon/pkg/domain"
	"readon/pkg/models"

	interfaces "readon/pkg/repository/interface"

	"gorm.io/gorm"
)

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{
		DB: DB,
	}
}

func (c adminDatabase) Login(ctx context.Context, admin models.Userlogindata) (int, bool) {
	var admins domain.Admin
	result := c.DB.Where("email = ? AND password = ?", admin.Email, admin.Password).Limit(1).Find(&admins)
	if result.Error != nil {
		fmt.Println("error while ckecking for mathing data :", result.Error)
	}

	if result.RowsAffected == 0 {
		return 0, false
	}
	return int(admins.ID), true
}

func (c *adminDatabase) ListUser(ctx context.Context) ([]models.ListOfUser, error) {
	var list []models.ListOfUser

	err := c.DB.Table("users").Select("id,name,email").Limit(8).Scan(&list).Error

	return list, err
}

func (c *adminDatabase) FindByID(ctx context.Context, id uint) (domain.User, error) {
	var user domain.User
	err := c.DB.First(&user, id).Error

	return user, err
}

func (c *adminDatabase) Delete(ctx context.Context, user domain.User) error {
	err := c.DB.Delete(&user).Error

	return err
}

func (c *adminDatabase) BlockOrUnBlock(ctx context.Context, id int) bool {
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
