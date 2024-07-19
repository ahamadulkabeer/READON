package usecase

import (
	"errors"
	"fmt"
	"log"
	"readon/pkg/api/helpers"
	"readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
	"time"

	"github.com/jinzhu/copier"
	"github.com/razorpay/razorpay-go"
)

type OrderUseCase struct {
	OrderRepo   interfaces.OrderRepository
	CartRepo    interfaces.CartRepository
	AddressRepo interfaces.AddressRepository
	ProductRepo interfaces.ProductRepository
	CouponRepo  interfaces.CouponRepository
}

func NewOrderUseCase(orepo interfaces.OrderRepository, crepo interfaces.CartRepository, arepo interfaces.AddressRepository, prepo interfaces.ProductRepository, coupRepo interfaces.CouponRepository) services.OrderUseCase {
	return &OrderUseCase{
		OrderRepo:   orepo,
		CartRepo:    crepo,
		AddressRepo: arepo,
		ProductRepo: prepo,
		CouponRepo:  coupRepo,
	}
}

func (c OrderUseCase) CreateOrder(userID, addressID, PaymentMethodID int, couponCodes []string) (string, error) {

	cart, err := c.CartRepo.GetItems(userID)
	if err != nil {
		fmt.Println("err : ", err)
		return "", err
	}

	if len(cart) == 0 {
		return "", errors.New("cart is empty")
	}

	totalQTY := 0
	totalPrice := 0.00

	for _, items := range cart {
		totalQTY += items.Quantity
		totalPrice += items.Price * float64(items.Quantity)
	}
	order := domain.Order{
		UserID:          uint(userID),
		AdressID:        uint(addressID),
		PaymentMethodID: uint(PaymentMethodID),
		TotalQuantity:   totalQTY,
		TotalPrice:      totalPrice,
		DeleveryCharge:  49.00, // fixed (for now)
		Status:          "processing",
	}

	helpers.CalculateCouponDiscount(c.CouponRepo, couponCodes, &cart, &order)

	fmt.Println("order :", order)

	if PaymentMethodID == 1 {
		order.PaymentStatus = "not paid"
	}
	var razorOrderId string
	if PaymentMethodID == 2 {
		if totalPrice > 1000 {
			return "", errors.New("OrderExceedsLimit: Order amount exceeds the maximum allowed limit of 1000")
		}
		razorOrderId, err = MakeRazorpayPayment(order.TotalPrice)
		if err != nil {
			return "", err
		}
		order.PaymentStatus = "payment pending"
		order.RazorPayOrderID = razorOrderId
	}

	// creates order and clears cart

	err = c.OrderRepo.CreateOrder(&order, cart)
	if err != nil {
		return "", errors.New("couldnt create order : db error")
	}
	for _, coupon := range couponCodes {
		fmt.Println("coupn :", coupon)
		c.CouponRepo.MarkCouponAsRedemed(coupon, order.ID)
	}
	return razorOrderId, nil

}

func (c OrderUseCase) RetryOrder(userID, orderID int) (string, error) {
	paymentStatus, err := c.OrderRepo.CheckPymentStatus(orderID)
	if err != nil {
		return "", err
	}
	fmt.Println("payment staus :", paymentStatus)
	if paymentStatus != "failed" {
		return "", errors.New("not a failed payment")
	}
	RazorOrderID, err := c.OrderRepo.FetchRazorOrderID(orderID)
	if err != nil {
		return "", err
	}

	return RazorOrderID, nil
}

func (c OrderUseCase) CancelOrder(userid, orderID int) error {
	// to do : add the paid amount to the wallet
	err := c.OrderRepo.CancelOrder(orderID, userid)
	if err != nil {
		return err
	}

	err = c.CouponRepo.MarkCouponAsNotRedeemed(uint(orderID))
	if err != nil {
		log.Println("order cancelled , failed to refresh used coupons !")
	}
	return nil
}

func (c OrderUseCase) ListOrders(userID int, pageDetails models.Pagination) ([]models.OrdersListing, error) {
	if pageDetails.NewPage == 0 {
		pageDetails.NewPage = 1
	}
	pageDetails.Size = 5
	pageDetails.Offset = pageDetails.Size * (pageDetails.NewPage - 1)
	fmt.Println("page details :", pageDetails)

	listOfOrders, err := c.OrderRepo.ListOrders(userID, pageDetails)
	if err != nil {
		return []models.OrdersListing{}, err
	}
	if len(listOfOrders) == 0 {
		return []models.OrdersListing{}, errors.New("no failed orders")
	}
	orderListing := make([]models.OrdersListing, len(listOfOrders))
	for ind := range listOfOrders {

		copier.Copy(&orderListing[ind], &listOfOrders[ind])

		orderItems, err := c.OrderRepo.ListOrderItems(int(listOfOrders[ind].ID))
		if err != nil {
			return orderListing, err
		}

		var orderItemsListing []models.OrderItemsListing

		for i := range orderItems {
			booklisting, err := c.ProductRepo.GetProduct(int(orderItems[i].BookID))
			if err != nil {
				return orderListing, err
			}
			orderItemsListing = append(orderItemsListing, models.OrderItemsListing{
				BookID:   int(orderItems[i].BookID),
				Title:    booklisting.Title,
				Price:    orderItems[i].Price,
				Quantity: orderItems[i].Quantity,
				Total:    orderItems[i].Price * float64(orderItems[i].Quantity),
			})
		}

		address, err := c.AddressRepo.GetAddress(listOfOrders[ind].AdressID, listOfOrders[ind].UserID)
		if err != nil {
			return orderListing, err
		}
		if listOfOrders[ind].PaymentMethodID == 1 {
			orderListing[ind].PaymentMethod = "COD"
		}
		if listOfOrders[ind].PaymentMethodID == 2 {
			orderListing[ind].PaymentMethod = "Online"
		}
		copier.Copy(&orderListing[ind].Address, address)

		orderListing[ind].Items = orderItemsListing

	}
	return orderListing, err

}

func (c OrderUseCase) GetOrder(userID, orderID int) (models.OrdersListing, error) {
	var orderListing models.OrdersListing
	order, err := c.OrderRepo.GetOrder(userID, orderID)
	if err != nil {
		return orderListing, err
	}
	copier.Copy(&orderListing, &order)
	orderItems, err := c.OrderRepo.ListOrderItems(orderID)
	if err != nil {
		return orderListing, err
	}
	var orderItemsListing []models.OrderItemsListing
	for i := range orderItems {
		booklisting, err := c.ProductRepo.GetProduct(int(orderItems[i].BookID))
		if err != nil {
			return orderListing, err
		}
		orderItemsListing = append(orderItemsListing, models.OrderItemsListing{
			BookID:   int(orderItems[i].BookID),
			Title:    booklisting.Title,
			Price:    orderItems[i].Price,
			Quantity: orderItems[i].Quantity,
			Total:    orderItems[i].Price * float64(orderItems[i].Quantity),
		})
	}
	address, err := c.AddressRepo.GetAddress(order.AdressID, order.UserID)
	if err != nil {
		return orderListing, err
	}
	if order.PaymentMethodID == 1 {
		orderListing.PaymentMethod = "COD"
	}
	if order.PaymentMethodID == 2 {
		orderListing.PaymentMethod = "Online"
	}
	copier.Copy(&orderListing.Address, address)
	orderListing.Items = orderItemsListing
	return orderListing, nil
}

// to do  : make a pagination here
func (c OrderUseCase) GetAllOrders(filter int) ([]domain.Order, error) {
	var list []domain.Order
	var err error
	list, err = c.OrderRepo.GetAllOrders(getTimeSpan(filter))
	if err != nil {
		return nil, err
	}
	fmt.Println("list :", list)
	return list, nil
}

var razorpayKey, razorpaySecret string

func LoadRazorpayConfig(key, secret string) error {
	if key == "" || secret == "" {
		return errors.New("the config data is empty")
	}
	razorpayKey = key
	razorpaySecret = secret
	return nil
}

func MakeRazorpayPayment(Price float64) (string, error) {

	client := razorpay.NewClient(razorpayKey, razorpaySecret)
	amountInPaise := Price * 100
	orderData := map[string]interface{}{
		"amount":          amountInPaise,
		"currency":        "INR",
		"receipt":         "order_rcptid_11",
		"notes":           map[string]string{"key": "value"}, // Additional notes if needed
		"payment_capture": 1,                                 // Auto-capture payment after order creation
	}

	// Create an order
	order, err := client.Order.Create(orderData, nil)
	if err != nil {
		return "", err
	}
	orderID := order["id"].(string)
	return orderID, nil

}

func (c OrderUseCase) VerifyPayment(paymentData models.PaymentVerificationData) error {

	if paymentData.PaymentStatus == "captured" {
		paymentData.PaymentStatus = "paid"
	}

	err := c.OrderRepo.UpdatePaymentStatus(paymentData)

	if err != nil {
		return err
	}

	return nil

}

func (c OrderUseCase) GetInvoiveData(userID, orderID int) (models.InvoiceData, error) {

	invoice := models.InvoiceData{
		CompanyName:    "ReadON",
		CompanyAddress: "123 Main Street, City, Country",
		CompanyContact: "+1 234 5678 910",
	}
	var err error
	invoice.Order, err = c.OrderRepo.GetOrder(userID, orderID)
	if err != nil {
		return invoice, err
	}
	orderItems, err := c.OrderRepo.ListOrderItems(orderID)
	if err != nil {
		return invoice, err
	}
	invoice.Date = invoice.Order.CreatedAt.Format("2006-01-02 15:04 ") + "+5:30"
	var orderItemsListing []models.OrderItemsListing
	for i := range orderItems {
		booklisting, err := c.ProductRepo.GetProduct(int(orderItems[i].BookID))
		if err != nil {
			return invoice, err
		}
		orderItemsListing = append(orderItemsListing, models.OrderItemsListing{
			BookID:   int(orderItems[i].BookID),
			Title:    booklisting.Title,
			Price:    orderItems[i].Price,
			Quantity: orderItems[i].Quantity,
			Total:    orderItems[i].Price * float64(orderItems[i].Quantity),
		})
	}
	invoice.OrderItems = orderItemsListing
	invoice.Address, err = c.AddressRepo.GetAddress(1, 1)
	if err != nil {
		return invoice, err
	}
	invoice.Total = invoice.Order.TotalPrice + invoice.Order.DeleveryCharge
	fmt.Println("invoive data ,", invoice)
	return invoice, nil
}

func (c OrderUseCase) GetChartData(pageDetails models.Pagination) (models.ChartData, error) {
	var chartData models.ChartData

	start, end := getTimeSpan(pageDetails.Filter)
	fmt.Println("start : ,", start, "end : ", end.Add(-time.Millisecond))

	data, err := c.OrderRepo.GetTotalAmountOfSpan(start, end.Add(-time.Millisecond), interval(pageDetails.Filter))
	if err != nil {
		fmt.Println("err", err)
		return chartData, err
	}
	fmt.Println(data)
	chartData.Span = []time.Time{start, end.Add(-time.Millisecond)}
	chartData.Data = data
	return chartData, nil
}

func getTimeSpan(filter int) (time.Time, time.Time) {

	fmt.Println("today now UTC :", time.Now().UTC())
	fmt.Println("today now  IST:", time.Now())

	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("err : ", err)
		return time.Time{}, time.Time{}
	}
	fmt.Println("location ", location)
	today := time.Now().In(location)

	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, location)

	fmt.Println("today truncated :", today)

	switch filter {
	case 1:
		todayStart := today
		todayEnd := todayStart.Add(24 * time.Hour)
		return todayStart, todayEnd
	case 2:
		thisWeekStart := today.AddDate(0, 0, -int(today.Weekday()))
		thisWeekEnd := thisWeekStart.AddDate(0, 0, 7)
		return thisWeekStart, thisWeekEnd
	case 3:
		thisMonthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, location)
		thisMonthEnd := thisMonthStart.AddDate(0, 1, 0)
		return thisMonthStart, thisMonthEnd
	case 4:
		thisYearStart := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, location)
		thisYearEnd := thisYearStart.AddDate(1, 0, 0)
		return thisYearStart, thisYearEnd
	}
	//baically no span
	return today, today
}

func interval(filter int) string {
	switch filter {
	case 1:
		return "1 hour"
	case 2, 3:
		return "1 day"
	case 4:
		return "1 month"
	}
	return ""
}

func (c OrderUseCase) GetTopTenCategory(filter models.Pagination) ([]models.TopTenCategory, error) {
	data, err := c.OrderRepo.FindTopTenCategories(filter)
	if err != nil {
		return []models.TopTenCategory{}, err
	}
	return data, nil
}

func (c OrderUseCase) GetTopTenBooks(filter models.Pagination) ([]models.ListingBook, error) {
	data, err := c.OrderRepo.FindTopTenBooks(filter)
	if err != nil {
		return []models.ListingBook{}, err
	}
	var list []models.ListingBook
	copier.Copy(&list, data)
	fmt.Println("data :", data)
	fmt.Println("lsit :", list)
	return list, nil
}
