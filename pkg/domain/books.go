package domain

type Book struct {
	ID         uint   `gorm:"primaryKey"`
	Title      string `gorm:"not null"`
	Author     string `gorm:"not null;default:'Untitled'"`
	About      string `gorm:"type:text"`
	Chapters   int
	Rating     float32
	Premium    bool     `gorm:"not null;default:false"`
	CategoryID int      `gorm:"not null;default:1" `
	Category   Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT;"`
}

//have to add  tags and reviews here on this table
