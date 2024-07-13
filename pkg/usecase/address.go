package usecase

import (
	"errors"
	"net/http"
	"readon/pkg/api/errorhandler"
	"readon/pkg/api/responses"
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

func (c AddressUsecase) AddAddress(address domain.Address) responses.Response {
	err := c.AddressRepo.CreateNewAddress(&address)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "Address not created", err.Error(), nil)
	}
	return responses.ClientReponse(http.StatusCreated, "address created successfully", nil, address)
}

func (c AddressUsecase) EditAddress(address domain.Address) responses.Response {
	ok, err := c.AddressRepo.AddressBelongsToUser(address.UserID, address.ID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "address not updated", err, nil)
	}
	if !ok {
		return responses.ClientReponse(http.StatusUnauthorized, "address not updated", err, nil)
	}
	err = c.AddressRepo.UpdateAddress(&address)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "address not updated", err, nil)
	}
	return responses.ClientReponse(http.StatusOK, "address updated", nil, address)
}

func (c AddressUsecase) ListAddress(userID uint) responses.Response {
	list, err := c.AddressRepo.ListAddresses(userID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't retrieve adresses", err, nil)
	}
	return responses.ClientReponse(http.StatusOK, "addresses fetched", nil, list)
}

func (c AddressUsecase) GetAddress(addressID, userID uint) responses.Response {
	address, err := c.AddressRepo.GetAddress(addressID, userID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "no address found", err, nil)

	}
	return responses.ClientReponse(http.StatusOK, "address fetched", nil, address)
}

func (c AddressUsecase) DeleteAddress(addressID, userID uint) responses.Response {
	ok, err := c.AddressRepo.AddressBelongsToUser(userID, addressID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "address not deleted", err, nil)
	}
	if !ok {
		return responses.ClientReponse(http.StatusNotFound, "address not deleted", errors.New("record not found"), nil)
	}
	err = c.AddressRepo.DeleteAddress(addressID, userID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "address not deleted", err, nil)
	}
	return responses.ClientReponse(http.StatusOK, "address deleted successfully", nil, nil)
}
