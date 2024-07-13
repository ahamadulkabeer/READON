package interfaces

import (
	"readon/pkg/api/responses"
	"readon/pkg/domain"
)

type AddressUsecase interface {
	AddAddress(address domain.Address) responses.Response
	EditAddress(address domain.Address) responses.Response
	ListAddress(userID uint) responses.Response
	GetAddress(addressID, userID uint) responses.Response
	DeleteAddress(addressID, userID uint) responses.Response
}
