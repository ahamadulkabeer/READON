package models

type Userlogindata struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type ListingBook struct {
	ID     int
	Title  string
	Author string
	Rating float32
}

type ListOfUser struct {
	id    uint
	Name  string
	Email string
}
type Newcategory struct {
	Name string `form:"name"`
}

type Product struct {
	Name       string `form:"name"`
	Author     string `form:"author"`
	Image      []byte `form:"image"`
	About      string `form:"about"`
	CategoryID int    `form:"category"`
}
