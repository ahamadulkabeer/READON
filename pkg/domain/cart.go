package domain

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID   uint `form:"userid"` // fk
	BookID   uint `form:"bookid"` // fk
	Quantity int
	Price    float64
	Book     Book `gorm:"forienkey:BookID;OnDelete:CASCADE,OnUpdate:CASCADE"`
	User     User `gorm:"forienkey:UserID;OnDelete:CASCADE,OnUpdate:CASCADE"`
}

// WANT TO MAKE TWO COLUMNS UNIQUE
// TOGETHER USERID & BOOKiD

//can be cascasded as the data is not needed anymore..
