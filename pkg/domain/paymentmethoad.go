package domain

type PaymentMethod struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"not null;unique"`
}
