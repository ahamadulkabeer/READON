package domain

type Book struct {
	ID         uint   `gorm:"primaryKey;column:id"`
	Title      string `gorm:"not null;default:'Untitled'"`
	Author     string `gorm:"not null"`
	About      string `gorm:"type:text"`
	Chapters   int
	Rating     float32
	Premium    bool `gorm:"not null;default:false"`
	Price      float64
	CategoryID int      `gorm:"not null;default:1" `
	Category   Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT;"`
}

//have to add  tags and reviews here on this table
