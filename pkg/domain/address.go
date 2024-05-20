package domain

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	UserID    uint
	Name      string `gorm:"not null" form:"name"`
	HouseNo   string `form:"houseNo"`
	HouseName string `form:"houseName"`
	Place     string `gorm:"not null" form:"place"`
	Landmark  string `form:"landmark"`
	City      string `gorm:"not null" form:"city"`
	District  string `gorm:"not null" form:"distric"`
	Country   string `gorm:"not null" form:"country"`
	Pincode   string `gorm:"not null" form:"pincode"`
	Mobile    string `gorm:"not null" form:"mobile"`
	User      User   `gorm:"forienkey:UserId;OnDelete:CASCADE,OnUpdate:CASCADE"`
}
