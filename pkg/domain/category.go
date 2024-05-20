package domain

type Category struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	// Description string
}
