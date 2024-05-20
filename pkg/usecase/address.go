package usecase

import (
	"errors"
	"fmt"
	"readon/pkg/api/errorhandler"
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

func (c AddressUsecase) AddAddress(address domain.Address) error {
	err := c.AddressRepo.CreateNewAdress(address)
	if err != nil {
		return errorhandler.ErrorData{
			Code:         0,
			Message:      "",
			ErrorMessage: err,
		}
	}
	return nil
}

func (c AddressUsecase) EditAddress(address domain.Address) error {
	ok, err := c.AddressRepo.AddressBelongsToUser(address.UserID, address.ID)
	if err != nil {
		return errorhandler.ErrorData{
			Code:         0,
			Message:      "internal server error", //?
			ErrorMessage: err,
		}
	}
	if !ok {
		return errorhandler.ErrorData{
			Code:         0,
			Message:      "address don't belong to user", //?
			ErrorMessage: err,
		}
	}
	err = c.AddressRepo.UpdateAddress(address)
	if err != nil {
		return errorhandler.ErrorData{
			Code:         0,
			Message:      "",
			ErrorMessage: err,
		}
	}
	return nil
}

func (c AddressUsecase) ListAddress(userID uint) ([]domain.Address, error) {
	list, err := c.AddressRepo.ListAdresses(userID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c AddressUsecase) GetAddress(addressID, userID uint) (domain.Address, error) {
	ok, err := c.AddressRepo.AddressBelongsToUser(userID, addressID)
	fmt.Println(ok, err, userID, addressID)
	if !ok || err != nil {
		return domain.Address{}, errors.New("no record found on address id and userid")
	}
	address, err := c.AddressRepo.GetAdress(addressID)
	if err != nil {
		return address, err
	}
	return address, err
}

func (c AddressUsecase) DeleteAddress(addressID, userID uint) error {
	ok, err := c.AddressRepo.AddressBelongsToUser(userID, addressID)
	if !ok || err != nil {
		return errors.New("no record found on address id and userid")
	}
	err = c.AddressRepo.DeleteAddress(addressID, userID)
	if err != nil {
		return errors.New("deletion failed : db error")
	}
	return nil
}
