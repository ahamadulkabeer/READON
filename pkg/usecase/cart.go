package usecase

import (
	"net/http"
	"readon/pkg/api/errorhandler"
	"readon/pkg/api/responses"
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

func (c CartUseCase) AddItem(userId, bookId int) responses.Response {

	count, err := c.CartRepo.CheckForItem(userId, bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't add item", err, nil)
	}

	if count != 0 {
		count++
		err = c.CartRepo.UpdateQty(userId, int(bookId), count)
		if err != nil {
			statusCode, _ := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't add item", err, nil)
		}
		return responses.ClientReponse(http.StatusOK, " item added ", err, nil)
	}

	price, err := c.ProductRepo.GetPrice(int(bookId))
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't add item", err, nil)
	}

	newCartItem := domain.Cart{
		UserID:   uint(userId),
		BookID:   uint(bookId),
		Price:    price,
		Quantity: 1,
	}

	err = c.CartRepo.AddItem(newCartItem)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't add item", err, nil)
	}
	return responses.ClientReponse(http.StatusOK, "item added to cart", nil, nil)
}

func (c CartUseCase) UpdateQty(userId, bookId, qty int) responses.Response {
	if qty <= 0 {
		err := c.CartRepo.DeleteItem(userId, bookId)
		if err != nil {
			statusCode, _ := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't update item quantity", err, nil)
		}
		return responses.ClientReponse(http.StatusOK, "item quantity updated", nil, nil)
	}

	count, err := c.CartRepo.CheckForItem(userId, bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't update item quantity", err, nil)
	}

	if count == 0 {
		return responses.ClientReponse(http.StatusNotFound, "couldn't update item quantity", "item not found in the cart", nil)
	}

	err = c.CartRepo.UpdateQty(userId, bookId, qty)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't update item quantity", err, nil)
	}

	return responses.ClientReponse(http.StatusOK, "item quantity updated", nil, nil)
}

func (c CartUseCase) DeleteItem(userId, bookId int) responses.Response {
	count, err := c.CartRepo.CheckForItem(userId, bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't remove item from cart", err, nil)
	}
	if count == 0 {
		return responses.ClientReponse(http.StatusNotFound, "couldn't remove item from cart", "item not found on cart !", nil)
	}
	err = c.CartRepo.DeleteItem(userId, bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't remove item from cart", err, nil)
	}
	return responses.ClientReponse(http.StatusOK, "item removed from the cart", nil, nil)
}

func (c CartUseCase) GetCart(userId int) responses.Response {
	var cart []domain.Cart
	cart, err := c.CartRepo.GetItems(userId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch user's cart", err, nil)
	}
	if len(cart) == 0 {
		return responses.ClientReponse(http.StatusOK, "cart is empty", nil, nil)
	}
	return responses.ClientReponse(http.StatusOK, "cart fetched successfully", nil, cart)
}
