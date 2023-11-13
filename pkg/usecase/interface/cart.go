package interfaces

import "readon/pkg/domain"

type CartUseCase interface {
	AddItem(item domain.Cart, userId int) error
	UpdateQty(userId, bookId, qty int) error
	DeleteItem(userId, bookId int) error
	GetCart(userId int) ([]domain.Cart, error)
}
