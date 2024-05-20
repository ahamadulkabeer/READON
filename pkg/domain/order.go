package domain

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID          uint // fk
	AdressID        uint // fk
	PaymentMethodID uint `gorm:"not null"` //  fk
	RazorPayOrderID string
	PaymentStatus   string
	PaymentID       string
	TotalQuantity   int
	TotalPrice      float64
	DeleveryCharge  float64
	Status          string
	User            User          `gorm:"forienkey:UserID;OnDelete:CASCADE,OnUpdate:CASCADE"`
	Adress          Address       `gorm:"forienkey:AdressID;OnDelete:CASCADE,OnUpdate:CASCADE"`
	PaymentMethod   PaymentMethod `gorm:"forienkey:PaymentMethodID;OnDelete:CASCADE,OnUpdate:CASCADE"`
}

// on fk
// on delete should not be cascaded
// as the order should be done
// on delevery address it should not be cascaded as adress conflict may occure
// same for others !

type OrderItems struct {
	ID       uint `gorm:"primarykey"`
	OrderID  uint //fk
	BookID   uint // fk
	Quantity int  `gorm:"not null;default:1"`
	Price    float64
	Book     Book  `gorm:"forienkey:BookID;OnDelete:CASCADE,OnUpdate:CASCADE"`
	Order    Order `gorm:"forienkey:OrderID;OnDelete:CASCADE,OnUpdate:CASCADE"`
}
