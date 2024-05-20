package models

import "time"

type SignupData struct {
	Name     string `copier:"must" form:"name"`
	Email    string `copier:"must" form:"email"`
	Password string `copier:"must" form:"password"`
}

type UserUpdateData struct {
	Id   int
	Name string `copier:"must" form:"name"`
}

type LoginData struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type Newcategory struct {
	Name string `form:"name"`
	//Description string `form:"description"`
}

type Product struct {
	Title      string `form:"title"`
	Author     string `form:"author"`
	Image      []byte `form:"image"`
	About      string `form:"about"`
	Price      float64
	CategoryID int `form:"category"`
}

type ProductUpdate struct {
	ID         int
	Name       string  `form:"name"`
	Author     string  `form:"author"`
	About      string  `form:"about"`
	Price      float64 `form:"price"`
	CategoryID int     `form:"category"`
}

type Coupon struct {
	Name           string
	CouponPrefix   string
	ValidFrom      time.Time
	ValidTill      time.Time
	MaxQuantity    int
	DiscountType   string
	DiscountAmount int
	IsBound        bool
}
