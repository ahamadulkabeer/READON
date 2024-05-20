package models

import (
	"readon/pkg/domain"
	"time"
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
	Order          domain.Order
	Date           string
	OrderItems     []OrderItemsListing
	GST            float64
	Total          float64
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
