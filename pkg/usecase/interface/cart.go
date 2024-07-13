package interfaces

import (
	"readon/pkg/api/responses"
)

type CartUseCase interface {
	AddItem(userID, bookID int) responses.Response
	UpdateQty(userID, bookID, qty int) responses.Response
	DeleteItem(userID, bookID int) responses.Response
	GetCart(userID int) responses.Response
}
