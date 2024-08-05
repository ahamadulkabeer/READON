package interfaces

import "readon/pkg/domain"

type AddressRepository interface {
	CreateNewAddress(address *domain.Address) error
	UpdateAddress(newAddress *domain.Address) error
	ListAddresses(userID uint) ([]domain.Address, error)
	GetAddress(addressID, userID uint) (domain.Address, error)
	DeleteAddress(addressID, userID uint) error
	AddressBelongsToUser(userID, addressID uint) (bool, error)
	GetNumberOfAdresses(userID uint) (int, error)
	AddressFound(addressID uint) (bool, error)
}
