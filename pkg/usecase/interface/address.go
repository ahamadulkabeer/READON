package interfaces

import (
	"readon/pkg/domain"
)

type AddressUsecase interface {
	AddAddress(address domain.Address, userId int) error
	EditAddress(adress domain.Address, userId int) error
	ListAddress(userId int) ([]domain.Address, error)
	GetAddress(addressId int) (domain.Address, error)
	DeleteAddress(addressId int) error
}
