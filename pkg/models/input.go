package models

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
	Name        string `form:"name"`
	Description string `form:"description"`
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
	Name               string  `form:"name"`
	Description        string  `form:"description"`
	CouponPrefix       string  `form:"prefix"`
	ValidFrom          string  `form:"validFrom"`
	ValidTill          string  `form:"validTill"`
	Limited            bool    `form:"limited"`
	MaxQuantity        int     `form:"maxQuantity"`
	MaxDiscount        float64 `form:"maxDiscount"`
	DiscountType       string  `form:"type"`
	DiscountAmount     int     `form:"amount"`
	ApplicableOn       string  `form:"applicableOn"`
	ApplicableCategory string  `form:"applicableCategory"`
	ApplicableProduct  string  `form:"applicableProduct"`
	IsBound            bool    `form:"isBound"`
}
type Address struct {
	Name      string `gorm:"not null" form:"name"`
	HouseNo   string `form:"houseNo"`
	HouseName string `form:"houseName"`
	Place     string `gorm:"not null" form:"place"`
	Landmark  string `form:"landmark"`
	City      string `gorm:"not null" form:"city"`
	District  string `gorm:"not null" form:"district"`
	Country   string `gorm:"not null" form:"country"`
	Pincode   string `gorm:"not null" form:"pincode"`
	Mobile    string `gorm:"not null" form:"mobile"`
}
