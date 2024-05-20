package repository

import (
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

// to do : clear this //
func (c adminDatabase) Login(admin models.LoginData) (int, bool) {
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

func (c *adminDatabase) ListAdmins() ([]models.Admin, error) {
	var list []models.Admin

	err := c.DB.Limit(8).Find(&list).Error

	return list, err
}
