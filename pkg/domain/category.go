package domain

type Category struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	CategoryID int `gorm:"not null;default:1" sql:"type:integer REFERENCES categories(id) ON UPDATE CASCADE ON DELETE SET DEFAULT"`
}
