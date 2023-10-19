package domain

type Book struct {
	ID         uint   `gorm:"primaryKey"`
	Title      string `gorm:"not null"`
	Author     string `gorm:"not null;default:'Untitled'"`
	About      string `gorm:"type:text"`
	Chapters   int
	Rating     float32
	CategoryID int      `gorm:"not null;default:1" sql:"type:integer REFERENCES categories(id) ON DELETE SET NULL"` // Foreign key without gorm tag
	Category   Category `gorm:"foreignkey:CategoryID"`
}

//have to add  tags and reviews here on this table
