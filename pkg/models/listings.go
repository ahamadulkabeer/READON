package models

import "readon/pkg/domain"

type ListingBook struct {
	ID         int
	Title      string
	Author     string
	Rating     float64
	About      string
	Price      float64
	Premium    bool
	CategoryId int
	Category   string
	Image      []byte
}

type User struct {
	ID         uint   `json:"id" copier:"must"`
	Name       string `json:"name" copier:"must"`
	Email      string `json:"email" copier:"must"`
	Permission bool   `json:"permission" copier:"must"`
}
type Admin struct {
	ID    uint   `json:"id" copier:"must"`
	Name  string `json:"name" copier:"must"`
	Email string `json:"email" copier:"must"`
}
type Pagination struct {
	Size     int    `json:"size" form:"size"`
	Filter   int    `json:"filter" form:"filter"`
	NewPage  int    `json:"page" form:"page"`
	Search   string `json:"search" form:"search"`
	Offset   int    `json:"offset" form:"offset"`
	Lastpage int    `json:"lastpage" form:"lastpage"`
}

type UserlistResponse struct {
	Pagination
	List []domain.User
}
type BooksListResponse struct {
	Pagination
	List []ListingBook
}

type ListCart struct {
	CartId   uint
	UserId   uint
	BookId   uint
	Quantity int
	Price    float64
}

type ListAddress struct {
	ID        uint
	Name      string
	HouseNo   string
	HouseName string
	Place     string
	Landmark  string
	City      string
	District  string
	Country   string
	Pincode   string
	Mobile    string
}

type OrdersListing struct {
	ID              uint `copier:"must"`
	TotalQuantity   int
	TotalPrice      float64
	Items           []OrderItemsListing
	PaymentMethodID uint
	PaymentMethod   string
	RazorPayOrderID string
	PaymentStatus   string
	PaymentID       string
	Address         ListAddress
	Status          string
}
type OrderItemsListing struct {
	BookID   int
	Title    string
	Price    float64
	Quantity int
	Total    float64
}
