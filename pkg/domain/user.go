package domain

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null" validate:"required,name"`
	Email    string `json:"email" gorm:"unique;not null" validate:"required"`
	Password string `json:"password" gorm:"not null" validate:"required,password"`
	Premium  bool

	Permission bool `gorm:"not null;default:true"`
}
