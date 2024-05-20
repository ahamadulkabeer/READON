package repository

import (
	"fmt"
	domain "readon/pkg/domain"
	interfaces "readon/pkg/repository/interface"

	"gorm.io/gorm"
)

type AddressDatabase struct {
	DB *gorm.DB
}

func NewAddressRepository(db *gorm.DB) interfaces.AddressRepository {
	return &AddressDatabase{
		DB: db,
	}
}

func (c AddressDatabase) CreateNewAdress(address domain.Address) error {
	result := c.DB.Create(&address)
	return result.Error
}

func (c AddressDatabase) UpdateAddress(newAddress domain.Address) error {
	result := c.DB.Save(&newAddress)
	return result.Error
}

func (c AddressDatabase) ListAdresses(userID uint) ([]domain.Address, error) {
	var list []domain.Address
	err := c.DB.Model(&domain.Address{}).Where("user_id = ?", userID).Find(&list).Error
	if err != nil {
		fmt.Println("err : ", err)
	}
	fmt.Println("lsit", list)
	return list, err
}

func (c AddressDatabase) GetAdress(addressID uint) (domain.Address, error) {
	var Address domain.Address
	err := c.DB.Model(&domain.Address{}).Where("id = ?", addressID).First(&Address).Error
	if err != nil {
		fmt.Println("err : ", err)
	}
	fmt.Println("add", Address)
	return Address, err
}
func (c AddressDatabase) DeleteAddress(addressID, userID uint) error {
	err := c.DB.Where("id = ? AND user_id = ?", addressID, userID).Delete(&domain.Address{}).Error
	if err != nil {
		fmt.Println("err : ", err)
	}
	return err
}

func (c AddressDatabase) AddressBelongsToUser(userID, addressID uint) (bool, error) {
	var count int64
	err := c.DB.Model(&domain.Address{}).Where("user_id = ? AND id = ?", userID, addressID).Count(&count).Error
	if err != nil {
		fmt.Println("err : ", err)
	}
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
