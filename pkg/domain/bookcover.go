package domain

type Bookcover struct {
	ID     uint   `gorm:"primarykey"`
	Image  []byte `gorm:"type:bytea" form:"image"`
	BookID uint   `gorm:"not null"`
	Book   Book   `gorm:"foreignkey:book_id"`
}
