package models

import (
	"readon/pkg/domain"
	"time"

	"github.com/shopspring/decimal"
)

type PaymentVerificationData struct {
	RazorOrderId   string
	RazorPaymentId string
	PaymentStatus  string
}

type ErrorResponse struct {
	Status string
	Err    string
	Hint   string
}

type InvoiceData struct {
	CompanyName    string
	CompanyAddress string
	CompanyContact string
	Address        domain.Address
	Order          ListOrders
	Date           string
	OrderItems     []ListOrderItems
	GST            float64
	Total          decimal.Decimal
}

type ChartData struct {
	Span []time.Time
	Data []SalesResult
}

type SalesResult struct {
	StartTime  time.Time
	TotalSales float64
}

type TopTenCategory struct {
	CategoryName string
	TotalSales   int
}

type UserDataError struct {
	UserNameErr string
	EmailErr    string
	PasswordErr string
	GeneralErr  string
}

type PaymentPageData struct {
	RazorpayOrderID string
	OrderID         uint
	UserName        string
	FinalPrice      float64
}
