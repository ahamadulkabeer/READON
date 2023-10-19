package models

import "readon/pkg/domain"

type SignupData struct {
	Name     string `copier:"must"`
	Email    string `copier:"must"`
	Password string `copier:"must"`
}

type Userlogindata struct {
	Email    string `form:"email"`
	Password string `form:"password"`
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

type ListingBook struct {
	ID     int
	Title  string
	Author string
	Rating float32
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

type ErrorResponse struct {
	Err    string
	Status string
	Hint   string
}

type Pagination struct {
	Size     int    `json:"size" form:"size"`
	Filter   int    `json:"filter" form:"filter"`
	NewPage  int    `json:"page" form:"page"`
	Search   string `json:"search" form:"search"`
	Lastpage int    `json:"lastpage" form:"lastpage"`
}

type UselistResponse struct {
	Pagination
	List []domain.User
}
type BooksListResponse struct {
	Pagination
	List []ListingBook
}
