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

func (c AddressDatabase) AddAdress(address domain.Address) error {
	result := c.DB.Save(&address)
	return result.Error
}

func (c AddressDatabase) EditAddress(newAddress domain.Address) error {
	result := c.DB.Save(&newAddress)
	// for now !!
	// creating a new record with all the values in the input
	return result.Error
}

func (c AddressDatabase) ListAdresses(userId int) ([]domain.Address, error) {
	var list []domain.Address
	err := c.DB.Model(&domain.Address{}).Table("addresses").Where("user_id = ?", userId).Find(&list).Error
	fmt.Println("list : ", list)
	return list, err
}

func (c AddressDatabase) GetAdress(addressId int) (domain.Address, error) {
	var Address domain.Address
	err := c.DB.Model(&domain.Address{}).Table("addresses").Where("address_id = ?", addressId).First(&Address).Error
	return Address, err
}
func (c AddressDatabase) DeleteAdress(addressId int) error {
	err := c.DB.Where("address_id = ?", addressId).Delete(&domain.Address{}).Error
	return err
}
