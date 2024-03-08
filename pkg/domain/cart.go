package domain

type Cart struct {
	//gorm.Model
	CartId   uint `gorm:"primarykey"`
	UserId   uint `form:"userid"` // fk
	BookId   uint `form:"bookid"` // fk
	Quantity int
	Price    float64
	Book     Book `gorm:"forienkey:BookId;OnDelete:CASCADE,OnUpdate:CASCADE"`
	User     User `gorm:"forienkey:UserId;OnDelete:CASCADE,OnUpdate:CASCADE"`
}

// WANT TO MAKE TWO COLUMNS UNIQUE
// TOGETHER USERID & BOOKiD

//can be cascasded as the data is not needed anymore..
