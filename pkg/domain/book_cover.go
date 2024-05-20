package domain

type Bookcover struct {
	ID     uint   `gorm:"primarykey"`
	Image  []byte `gorm:"type:bytea" form:"image"`
	BookID uint   `gorm:"not null;index"`
	Book   Book   `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
