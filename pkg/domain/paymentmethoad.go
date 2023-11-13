package domain

type PaymentMethoad struct {
	MethodId uint   `gorm:"primarykey"`
	Name     string `gorm:"not null;unique"`
}
