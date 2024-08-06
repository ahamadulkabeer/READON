package models

import (
	"readon/pkg/domain"
	"time"
)

type Pagination struct {
	Filter          string `json:"filter" form:"filter"`
	Search          string `json:"search" form:"search"`
	Size            int    `json:"size" form:"size"`
	Page            int    `json:"page" form:"page"`
	NumberOfResults int    `json:"numberOfResults" form:"numberOfResults"`
}

type ListBook struct {
	ID       int
	Title    string
	Author   string
	Rating   float64
	About    string
	Price    float64
	Premium  bool
	Category ListCategories
	Image    []byte
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

type UserlistResponse struct {
	Pagination
	List []domain.User
}
type BooksListResponse struct {
	Pagination
	List []ListBook
}

type ListCartItem struct {
	ID         int
	BookId     uint
	Quantity   int
	Price      float64
	TotalPrice float64
	Book       ListBook
}

type ListCart struct {
	TotalQuantity int
	TotalPrice    float64
	Items         []ListCartItem
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

type ListOrders struct {
	ID              uint `copier:"must"`
	TotalQuantity   int
	TotalPrice      float64
	DiscountedPrice float64
	TotalDiscount   float64
	DeleveryCharge  float64
	Items           []ListOrderItems
	PaymentMethodID uint
	PaymentMethod   string
	RazorPayOrderID string
	PaymentStatus   string
	PaymentID       string
	AdressID        uint
	Address         ListAddress
	Status          string
}
type ListOrderItems struct {
	BookID   int
	Price    float64
	Quantity int
	Total    float64
	Book     ListBook
}

type ListCategories struct {
	ID   int
	Name string
}

type ListUserCoupons struct {
	CouponCode string
	Redeemed   bool
	RedeemedOn uint
	Coupon     ListCoupons
}

type ListCoupons struct {
	ID                 uint
	Name               string
	Description        string
	Prefix             string
	DiscountType       string
	DiscountAmount     int
	ApplicableOn       string
	ApplicableCategory string
	ApplicableProduct  string
	ValidFrom          time.Time
	ValidTill          time.Time
	Limited            bool
	MaxQuantity        int
	IsBound            bool
	Expired            bool
}

// for swagger

type PaginatedListCoupons struct {
	Coupons    []ListCoupons
	Pagination Pagination
}

type PaginatedListBooks struct {
	Books      []ListBook
	Pagination Pagination
}

type PaginatedListOrders struct {
	Orders     []ListOrders
	Pagination Pagination
}

// for now

type ListBookCover struct {
	ID    uint
	Image []byte
}
