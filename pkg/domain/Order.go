package domain

type Order struct {
	OrderId          uint `gorm:"primarykey"`
	UserId           uint // fk
	BookId           uint // fk
	Quantity         int  `gorm:"not null;default:1"`
	PaymentMethoadId uint `gorm:"not null"`
	AdressId         uint // fk
	DeliveryStatus   bool
	TotalPrice       float64
	User             User           `gorm:"forienkey:UserId;OnDelete:CASCADE,OnUpdate:CASCADE"`
	Book             Book           `gorm:"forienkey:BookId;OnDelete:CASCADE,OnUpdate:CASCADE"`
	Adress           Adress         `gorm:"forienkey:AdressId;OnDelete:CASCADE,OnUpdate:CASCADE"`
	PaymentMethoad   PaymentMethoad `gorm:"forienkey:PaymentMethoadId;OnDelete:CASCADE,OnUpdate:CASCADE"`
}

// on fk
// on delete should not be cascaded
// as the orderr should be done
// on delevery address it should not be cascaded as adress conflict may occure
// same for others !
