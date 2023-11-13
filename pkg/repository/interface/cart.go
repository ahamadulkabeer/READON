package interfaces

import "readon/pkg/domain"

type CartRepository interface {
	AddItem(item domain.Cart, userId int) error
	UpdateQty(userId, bookId, NQty int) error
	DeleteItem(userId, bookId int) error
	GetItems(userId int) ([]domain.Cart, error)
	// ClearCart
	CheckForItem(userId, bookId int) (int, error)
}

// In the GetItems instead of domain.cart must use another on somethin useful like costomizedobject with everything need ed
// or you can just go with this and get oth er need effddata in the usecase layeer :]
