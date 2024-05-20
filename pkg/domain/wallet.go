package domain

import "gorm.io/gorm"

type WalletHistory struct {
	gorm.Model
	UserID           uint
	TansactionAmount int
	TransactionType  string
	Transaction      string
	User             User
}
