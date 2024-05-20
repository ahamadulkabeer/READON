package interfaces

import "readon/pkg/domain"

type AddressRepository interface {
	CreateNewAdress(address domain.Address) error
	UpdateAddress(newAddress domain.Address) error
	ListAdresses(userID uint) ([]domain.Address, error)
	GetAdress(addressID uint) (domain.Address, error)
	DeleteAddress(addressID, userID uint) error
	AddressBelongsToUser(userID, addressID uint) (bool, error)
}
