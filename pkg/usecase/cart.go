package usecase

import (
	"errors"
	domain "readon/pkg/domain"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
)

type CartUseCase struct {
	CartRepo interfaces.CartRepository
}

func NewCartUseCase(repo interfaces.CartRepository) services.CartUseCase {
	return &CartUseCase{
		CartRepo: repo,
	}
}

func (c CartUseCase) AddItem(item domain.Cart, userId int) error {
	count, err := c.CartRepo.CheckForItem(userId, int(item.BookId))
	if err != nil {
		return err
	}

	if count != 0 {
		count++
		err = c.CartRepo.UpdateQty(userId, int(item.BookId), count)
		if err != nil {
			return err
		}
		return nil
	}
	err = c.CartRepo.AddItem(item, userId)
	if err != nil {
		return err
	}
	return nil
}

func (c CartUseCase) UpdateQty(userId, bookId, qty int) error {
	// i dont know if the count == 0 thruws an error
	count, err := c.CartRepo.CheckForItem(userId, bookId)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("item not found on cart !")
	}
	err = c.CartRepo.UpdateQty(userId, bookId, qty)
	if err != nil {
		return err
	}
	return nil
}

func (c CartUseCase) DeleteItem(userId, bookId int) error {
	count, err := c.CartRepo.CheckForItem(userId, bookId)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("item not found on cart !")
	}
	err = c.CartRepo.DeleteItem(userId, bookId)
	if err != nil {
		return err
	}
	return nil
}

func (c CartUseCase) GetCart(userId int) ([]domain.Cart, error) {
	var list []domain.Cart
	list, err := c.CartRepo.GetItems(userId)
	if err != nil {
		return nil, err
	}
	return list, nil
}
