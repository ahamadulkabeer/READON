package repository

import (
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

func (c AddressDatabase) CreateNewAddress(address *domain.Address) error {
	err := c.DB.Create(address).Error
	if err != nil {
		return err
	}
	return nil
}

func (c AddressDatabase) UpdateAddress(newAddress *domain.Address) error {
	err := c.DB.Save(newAddress).Error
	if err != nil {
		return err
	}
	return nil
}

func (c AddressDatabase) ListAddresses(userID uint) ([]domain.Address, error) {
	var list []domain.Address
	err := c.DB.Model(&domain.Address{}).Where("user_id = ?", userID).Find(&list).Error
	if err != nil {
		return []domain.Address{}, err
	}
	return list, err
}

func (c AddressDatabase) GetAddress(addressID, userID uint) (domain.Address, error) {
	var Address domain.Address
	err := c.DB.Model(&domain.Address{}).Where("id = ? AND user_id = ?", addressID, userID).First(&Address).Error
	if err != nil {
		return domain.Address{}, err
	}
	return Address, err
}

func (c AddressDatabase) DeleteAddress(addressID, userID uint) error {
	err := c.DB.Where("id = ? AND user_id = ?", addressID, userID).Delete(&domain.Address{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c AddressDatabase) AddressBelongsToUser(userID, addressID uint) (bool, error) {
	var count int64
	err := c.DB.Model(&domain.Address{}).Where("user_id = ? AND id = ?", userID, addressID).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (c AddressDatabase) GetNumberOfAdresses(userID uint) (int, error) {
	var count int64
	err := c.DB.Model(&domain.Address{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (c AddressDatabase) AddressFound(addressID uint) (bool, error) {
	var count int64
	err := c.DB.Model(&domain.Address{}).Where("id = ?", addressID).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
