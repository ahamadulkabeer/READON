package interfaces

import (
	"readon/pkg/domain"
)

type AddressUsecase interface {
	AddAddress(address domain.Address) error
	EditAddress(address domain.Address) error
	ListAddress(userID uint) ([]domain.Address, error)
	GetAddress(addressID, userID uint) (domain.Address, error)
	DeleteAddress(addressID, userID uint) error
}
