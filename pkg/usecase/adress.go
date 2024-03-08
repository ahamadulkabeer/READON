package usecase

import (
	domain "readon/pkg/domain"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
)

type AddressUsecase struct {
	AddressRepo interfaces.AddressRepository
}

func NewAddressUsecase(repo interfaces.AddressRepository) services.AddressUsecase {
	return &AddressUsecase{
		AddressRepo: repo,
	}
}

func (c AddressUsecase) AddAddress(address domain.Address, userId int) error {
	address.UserId = uint(userId)
	err := c.AddressRepo.AddAdress(address)
	return err
}
func (c AddressUsecase) EditAddress(adress domain.Address, userId int) error {
	adress.UserId = uint(userId)
	err := c.AddressRepo.EditAddress(adress)
	return err
}
func (c AddressUsecase) ListAddress(userId int) ([]domain.Address, error) {
	list, err := c.AddressRepo.ListAdresses(userId)
	if err != nil {
		return nil, err
	}
	return list, nil
}
func (c AddressUsecase) GetAddress(addressId int) (domain.Address, error) {
	address, err := c.AddressRepo.GetAdress(addressId)
	if err != nil {
		return address, err
	}
	// may want to check if theadress is really of the users by comapring id from contect and adress
	return address, err
}
func (c AddressUsecase) DeleteAddress(addressId int) error {
	err := c.AddressRepo.DeleteAdress(addressId)
	return err
}
