package usecase

import (
	"net/http"
	"readon/pkg/api/errorhandler"
	"readon/pkg/api/responses"
	domain "readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"

	"github.com/jinzhu/copier"
)

type AddressUsecase struct {
	AddressRepo interfaces.AddressRepository
}

func NewAddressUsecase(repo interfaces.AddressRepository) services.AddressUsecase {
	return &AddressUsecase{
		AddressRepo: repo,
	}
}

func (c AddressUsecase) AddAddress(newAddress domain.Address) responses.Response {
	// check if the limit of allowed  number of addresses
	count, err := c.AddressRepo.GetNumberOfAdresses(newAddress.UserID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "Address not created", err.Error(), nil)
	}
	if count >= 5 {
		return responses.ClientReponse(http.StatusForbidden, "Address can't be created",
			"user only allowed to have upto 5 Addresses ", nil)
	}

	// creates new address
	err = c.AddressRepo.CreateNewAddress(&newAddress)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "Address not created", err.Error(), nil)
	}

	//  response with only the nessessary data
	var address models.ListAddress
	copier.Copy(&address, &newAddress)
	return responses.ClientReponse(http.StatusCreated, "address created successfully", nil, address)
}

func (c AddressUsecase) EditAddress(newAddress domain.Address) responses.Response {
	// checks if the address exist
	ok, err := c.AddressRepo.AddressFound(newAddress.ID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "address not updated", err.Error(), nil)
	}
	if !ok {
		return responses.ClientReponse(http.StatusNotFound, "address not updated",
			"address with given id doesn't exist", nil)
	}

	// checks if address belong to the user
	ok, err = c.AddressRepo.AddressBelongsToUser(newAddress.UserID, newAddress.ID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "address not updated", err.Error(), nil)
	}
	if !ok {
		return responses.ClientReponse(http.StatusUnauthorized, "address can't be updated",
			"address doesn't belong to user", nil)
	}
	// updates the address
	err = c.AddressRepo.UpdateAddress(&newAddress)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "address not updated", err.Error(), nil)
	}

	//  response with only the nessessary data
	var address models.ListAddress
	copier.Copy(&address, &newAddress)
	return responses.ClientReponse(http.StatusOK, "address updated", nil, address)
}

func (c AddressUsecase) ListAddress(userID uint) responses.Response {
	// retrives list of addresses
	list, err := c.AddressRepo.ListAddresses(userID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't retrieve adresses", err.Error(), nil)
	}

	// response with only the nessessary data
	listOfAddresses := []models.ListAddress{}
	copier.Copy(&listOfAddresses, &list)
	return responses.ClientReponse(http.StatusOK, "addresses fetched", nil, listOfAddresses)
}

func (c AddressUsecase) GetAddress(addressID, userID uint) responses.Response {
	// checks if the address exist
	ok, err := c.AddressRepo.AddressFound(addressID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, " address not found", err.Error(), nil)
	}
	if !ok {
		return responses.ClientReponse(http.StatusNotFound, "address not  found",
			"address with given id doesn't exist", nil)
	}

	// checks if address belong to the user
	ok, err = c.AddressRepo.AddressBelongsToUser(userID, addressID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "address not fetched", err.Error(), nil)
	}
	if !ok {
		return responses.ClientReponse(http.StatusUnauthorized, "address not fetched",
			"address doesn't belong to user", nil)
	}

	// retrieves the resource
	data, err := c.AddressRepo.GetAddress(addressID, userID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, " address not found", err.Error(), nil)

	}
	var address models.ListAddress
	copier.Copy(&address, &data)
	return responses.ClientReponse(http.StatusOK, "address fetched", nil, address)
}

func (c AddressUsecase) DeleteAddress(addressID, userID uint) responses.Response {
	// checks if the address exist
	ok, err := c.AddressRepo.AddressFound(addressID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "address not deleted", err.Error(), nil)
	}
	if !ok {
		return responses.ClientReponse(http.StatusNotFound, "address not deleted",
			"address with given id doesn't exist", nil)
	}

	// checks if address belong to the user
	ok, err = c.AddressRepo.AddressBelongsToUser(userID, addressID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "address not deleted", err.Error(), nil)
	}
	if !ok {
		return responses.ClientReponse(http.StatusUnauthorized, "address not deleted",
			"address doesn't belong to user", nil)
	}

	// deletes the address
	err = c.AddressRepo.DeleteAddress(addressID, userID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "address not deleted", err.Error(), nil)
	}

	// response with the id of deleted resource
	return responses.ClientReponse(http.StatusOK, "address deleted successfully", nil,
		map[string]interface{}{
			"address id : ": addressID,
		})
}
