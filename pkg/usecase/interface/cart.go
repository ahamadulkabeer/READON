package interfaces

import "readon/pkg/domain"

type CartUseCase interface {
	AddItem(userID, bookID int) error
	UpdateQty(userID, bookID, qty int) error
	DeleteItem(userID, bookID int) error
	GetCart(userID int) ([]domain.Cart, error)
}
