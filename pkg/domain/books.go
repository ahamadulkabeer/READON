package domain

type Books struct {
	ID         uint   `gorm:"primaryKey"`
	Title      string `gorm:"not null;"`
	Author     string `gorm:"not null;default:'Untitled'"`
	Chapters   int
	Rating     float32
	CategoryID int
}

//have to add  tags and reviews here on this table
