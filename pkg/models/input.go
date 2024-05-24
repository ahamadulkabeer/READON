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
	Name           string `form:"name"`
	Description    string `form:"description"`
	CouponPrefix   string `form:"prefix"`
	ValidFrom      string `form:"validfrom"`
	ValidTill      string `form:"validtill"`
	MaxQuantity    int    `form:"quantity"`
	DiscountType   string `form:"type"`
	DiscountAmount int    `form:"amount"`
	IsBound        bool   `form:"isbound"`
}
