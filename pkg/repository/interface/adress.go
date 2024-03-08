package interfaces

import "readon/pkg/domain"

type AddressRepository interface {
	AddAdress(address domain.Address) error
	EditAddress(newAddress domain.Address) error
	ListAdresses(userId int) ([]domain.Address, error)
	GetAdress(addressId int) (domain.Address, error)
	DeleteAdress(addressId int) error
}
