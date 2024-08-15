package usecase

import (
	"fmt"
	"net/http"
	"readon/pkg/api/errorhandler"
	"readon/pkg/api/responses"
	domain "readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"

	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
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

	// check if the item exist in the user's cart and gets the quantity
	quantity, err := c.CartRepo.GetItemQuantity(userId, bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't add item", err.Error(), nil)
	}

	// if !exist : add item
	if quantity <= 0 {
		// get the items price
		price, err := c.ProductRepo.GetPrice(int(bookId))
		if err != nil {
			statusCode, _ := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't add item", err.Error(), nil)
		}

		// initilising the cart item
		newCartItem := domain.Cart{
			UserID:   uint(userId),
			BookID:   uint(bookId),
			Price:    price,
			Quantity: 1,
		}

		// add item to cart
		err = c.CartRepo.AddItem(newCartItem)
		if err != nil {
			statusCode, _ := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't add item", err.Error(), nil)
		}
		// continue...
	}

	// if exist :  increate the quatity by 1
	if quantity > 0 {
		quantity++
		err = c.CartRepo.UpdateQty(userId, int(bookId), quantity)
		if err != nil {
			statusCode, _ := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't add item", err.Error(), nil)
		}
		//continue...
	}

	// get cart item
	newCartItem, err := c.CartRepo.GetItem(uint(bookId))
	if err != nil {
		_, _ = errorhandler.HandleDatabaseError(err)
		fmt.Println("err :", err.Error())
		//return responses.ClientReponse(statusCode, "couldn't add item", err.Error(), nil)
	}

	//response with nessessary data
	var cartItem models.ListCartItem
	copier.Copy(&cartItem, &newCartItem)
	cartItem.TotalPrice = decimal.NewFromFloat(newCartItem.Price * float64(cartItem.Quantity)).Round(2)
	return responses.ClientReponse(http.StatusOK, "item added to cart", nil, cartItem)
}

func (c CartUseCase) UpdateQty(userId, bookId, newQuantity int) responses.Response {

	// if new quantity is less than 1 : removeitem from cart
	if newQuantity <= 0 {
		err := c.CartRepo.DeleteItem(userId, bookId)
		if err != nil {
			statusCode, _ := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't update item quantity", err.Error(), nil)
		}
		return responses.ClientReponse(http.StatusOK, "item quantity updated", nil, nil)
	}

	// check for the item in the cart and get item quantity
	qty, err := c.CartRepo.GetItemQuantity(userId, bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't update item quantity", err.Error(), nil)
	}
	if qty == 0 {
		return responses.ClientReponse(http.StatusNotFound, "couldn't update item quantity", "item not found in the cart", nil)
	}

	// updates item quantity with new quantity
	err = c.CartRepo.UpdateQty(userId, bookId, newQuantity)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't update item quantity", err.Error(), nil)
	}

	// get cart item
	newCartItem, err := c.CartRepo.GetItem(uint(bookId))
	if err != nil {
		_, _ = errorhandler.HandleDatabaseError(err)
		fmt.Println("err :", err.Error())
		//return responses.ClientReponse(statusCode, "couldn't add item", err.Error(), nil)
	}

	//response with nessessary data
	var cartItem models.ListCartItem
	copier.Copy(&cartItem, &newCartItem)
	cartItem.TotalPrice = decimal.NewFromFloat(newCartItem.Price * float64(cartItem.Quantity)).Round(2)
	return responses.ClientReponse(http.StatusOK, "item quantity updated", nil, cartItem)
}

func (c CartUseCase) DeleteItem(userId, bookId int) responses.Response {

	//  get quantity and check if item exist in user's cart  get quantity
	qty, err := c.CartRepo.GetItemQuantity(userId, bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't remove item from cart", err.Error(), nil)
	}
	if qty == 0 {
		return responses.ClientReponse(http.StatusNotFound, "couldn't remove item from cart", "item not found on cart !", nil)
	}

	// delete cart
	err = c.CartRepo.DeleteItem(userId, bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't remove item from cart", err.Error(), nil)
	}

	//response with nessessary data
	cartData, err := c.CartRepo.GetItems(userId)
	if err != nil {
		_, _ = errorhandler.HandleDatabaseError(err)
	}
	var cart models.ListCart
	for _, x := range cartData {
		var cartItem models.ListCartItem
		copier.Copy(&cartItem, &x)
		cartItem.TotalPrice = decimal.NewFromFloat(x.Price * float64(x.Quantity))
		cart.TotalPrice = cart.TotalPrice.Add(cartItem.TotalPrice)

		cart.TotalQuantity += cartItem.Quantity
		cart.Items = append(cart.Items, cartItem)
	}
	return responses.ClientReponse(http.StatusOK, "item removed from the cart", nil, cart)
}

func (c CartUseCase) GetCart(userId int) responses.Response {
	// retrieve user's cart
	cartData, err := c.CartRepo.GetItems(userId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch user's cart", err.Error(), nil)
	}

	// if cart is empty
	if len(cartData) == 0 {
		return responses.ClientReponse(http.StatusOK, "cart is empty", nil, nil)
	}

	//response with nessessary data
	var cart models.ListCart
	for _, x := range cartData {
		var cartItem models.ListCartItem
		copier.Copy(&cartItem, &x)
		cartItem.TotalPrice = decimal.NewFromFloat(x.Price * float64(x.Quantity))
		cart.TotalPrice = cart.TotalPrice.Add(cartItem.TotalPrice)
		cart.TotalQuantity += cartItem.Quantity
		cart.Items = append(cart.Items, cartItem)
	}

	cart.TotalPrice = cart.TotalPrice.Add(decimal.NewFromFloat(1)).Round(2)
	return responses.ClientReponse(http.StatusOK, "cart fetched successfully", nil, cart)
}
