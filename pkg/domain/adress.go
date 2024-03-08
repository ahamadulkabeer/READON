package domain

type Address struct {
	AddressId uint `gorm:"primarykey"`
	UserId    uint
	Name      string `gorm:"not null"`
	HouseNo   string
	HouseName string
	Place     string `gorm:"not null"`
	Landmark  string
	City      string `gorm:"not null"`
	District  string `gorm:"not null"`
	Country   string `gorm:"not null"`
	Pincode   string `gorm:"not null"`
	Mobile    string `gorm:"not null"`
	User      User   `gorm:"forienkey:UserId;OnDelete:CASCADE,OnUpdate:CASCADE"`
}
