package usecase

import (
	"errors"
	domain "readon/pkg/domain"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
)

type CartUseCase struct {
	CartRepo    interfaces.CartRepository
	ProductRepo interfaces.ProductRepository
}

func NewCartUseCase(crepo interfaces.CartRepository, prepo interfaces.ProductRepository) services.CartUseCase {
	return &CartUseCase{
		CartRepo:    crepo,
		ProductRepo: prepo,
	}
}

func (c CartUseCase) AddItem(userId, bookId int) error {
	count, err := c.CartRepo.CheckForItem(userId, bookId)
	if err != nil {
		return err
	}
	if count != 0 {
		count++
		err = c.CartRepo.UpdateQty(userId, int(bookId), count)
		if err != nil {
			return err
		}
		return nil
	}
	price, err := c.ProductRepo.GetPrice(int(bookId))
	if err != nil {
		return err
	}
	newCartItem := domain.Cart{
		UserID:   uint(userId),
		BookID:   uint(bookId),
		Price:    price,
		Quantity: 1,
	}
	err = c.CartRepo.AddItem(newCartItem)
	if err != nil {
		return err
	}
	return nil
}

func (c CartUseCase) UpdateQty(userId, bookId, qty int) error {
	if qty <= 0 {
		err := c.CartRepo.DeleteItem(userId, bookId)
		if err != nil {
			return err
		}
	}
	count, err := c.CartRepo.CheckForItem(userId, bookId)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("item not found in cart")
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
	if len(list) == 0 {
		return []domain.Cart{}, errors.New("cart is empty")
	}
	return list, nil
}
