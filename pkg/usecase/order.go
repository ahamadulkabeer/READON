package usecase

import (
	"errors"
	"fmt"
	"net/http"
	"readon/pkg/api/errorhandler"
	"readon/pkg/api/helpers"
	"readon/pkg/api/responses"
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

func NewOrderUseCase(orepo interfaces.OrderRepository,
	crepo interfaces.CartRepository,
	arepo interfaces.AddressRepository,
	prepo interfaces.ProductRepository,
	coupRepo interfaces.CouponRepository) services.OrderUseCase {
	return &OrderUseCase{
		OrderRepo:   orepo,
		CartRepo:    crepo,
		AddressRepo: arepo,
		ProductRepo: prepo,
		CouponRepo:  coupRepo,
	}
}

func (c OrderUseCase) CreateOrder(userID, addressID, PaymentMethodID int, couponCodes []string) responses.Response {

	cart, err := c.CartRepo.GetItems(userID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couln't place order", err.Error(), nil)
	}

	if len(cart) == 0 {
		return responses.ClientReponse(http.StatusNotFound, "couln't place order", "cart is empty", nil)
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
		PaymentStatus:   "not paid",
		DeleveryCharge:  49.00, // fixed (for now)
		Status:          "processing",
	}

	helpers.CalculateCouponDiscount(c.CouponRepo, couponCodes, &cart, &order)

	fmt.Println("order :", order)

	var razorOrderID string
	if PaymentMethodID == 2 {
		if totalPrice > 1000 {
			return responses.ClientReponse(http.StatusUnprocessableEntity, "couln't place order",
				errors.New("OrderExceedsLimit: Order amount exceeds the maximum allowed limit of 1000"), nil)
		}
		razorOrderID, err = MakeRazorpayPayment(order.TotalPrice)
		if err != nil {
			return responses.ClientReponse(http.StatusInternalServerError, "couln't place order", errors.New("payment gatway error").Error(), nil)
		}
		order.PaymentStatus = "payment pending"
		order.RazorPayOrderID = razorOrderID
	}

	// creates order and clears cart

	err = c.OrderRepo.CreateOrder(&order, cart)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couln't place order", err.Error(), nil)
	}
	for _, coupon := range couponCodes {
		fmt.Println("coupn :", coupon)
		c.CouponRepo.MarkCouponAsRedemed(coupon, order.ID)
	}
	res := c.MakeOrderResponse(order)
	return responses.ClientReponse(http.StatusOK, "order placed , "+"razorpay order id : "+razorOrderID, nil, res.Data)

}

func (c OrderUseCase) RetryOrder(userID, orderID int) responses.Response {
	paymentStatus, err := c.OrderRepo.CheckPymentStatus(orderID)

	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't retry order", err.Error(), nil)
	}

	if paymentStatus != "failed" {
		return responses.ClientReponse(http.StatusBadRequest, "couldn't retry order", "not a failed payment", nil)
	}

	razorOrderID, err := c.OrderRepo.FetchRazorOrderID(orderID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't retry order", err.Error(), nil)
	}

	return responses.ClientReponse(http.StatusOK, "order placed",
		nil, "razorpay order id : "+razorOrderID)
}

func (c OrderUseCase) CancelOrder(userid, orderID int) responses.Response {
	// to do : add the paid amount to the wallet
	err := c.OrderRepo.CancelOrder(orderID, userid)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "coulnd't cancel order", err.Error(), nil)
	}

	err = c.CouponRepo.MarkCouponAsNotRedeemed(uint(orderID))
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "order cancelled , failed to refresh used coupons !", err.Error(), nil)
	}
	return responses.ClientReponse(http.StatusOK, "order cancelled", nil, nil)
}

func (c OrderUseCase) ListOrders(userID int, pageDetails models.Pagination) responses.Response {
	if pageDetails.NewPage == 0 {
		pageDetails.NewPage = 1
	}
	pageDetails.Size = 5
	pageDetails.Offset = pageDetails.Size * (pageDetails.NewPage - 1)
	fmt.Println("page details :", pageDetails)

	listOfOrders, err := c.OrderRepo.ListOrders(userID, pageDetails)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch orders", err.Error(), nil)
	}

	if len(listOfOrders) == 0 {
		return responses.ClientReponse(http.StatusNotFound, "no orders found", nil, nil)
	}

	orderListing := make([]models.OrdersListing, len(listOfOrders))
	for ind := range listOfOrders {

		copier.Copy(&orderListing[ind], &listOfOrders[ind])

		orderItems, err := c.OrderRepo.ListOrderItems(int(listOfOrders[ind].ID))
		if err != nil {
			statusCode, _ := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't fetch orders", err.Error(), nil)
		}

		var orderItemsListing []models.OrderItemsListing

		for i := range orderItems {
			booklisting, err := c.ProductRepo.GetProduct(int(orderItems[i].BookID))
			if err != nil {
				statusCode, _ := errorhandler.HandleDatabaseError(err)
				return responses.ClientReponse(statusCode, "couldn't fetch orders", err.Error(), nil)
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
			statusCode, _ := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't fetch orders", err.Error(), nil)
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
	return responses.ClientReponse(http.StatusOK, "orders fetched successfully", nil, orderListing)

}

func (c OrderUseCase) GetOrder(userID, orderID int) responses.Response {

	order, err := c.OrderRepo.GetOrder(userID, orderID)
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch order", err.Error(), nil)
	}
	return c.MakeOrderResponse(order)
}

// to do  : make a pagination here
func (c OrderUseCase) GetAllOrders(filter int) responses.Response {
	var list []domain.Order
	var err error
	list, err = c.OrderRepo.GetAllOrders(getTimeSpan(filter))
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch all orders", err.Error(), nil)
	}
	return responses.ClientReponse(http.StatusOK, "fetched all orders", nil, list)
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

func (c OrderUseCase) VerifyPayment(paymentData models.PaymentVerificationData) responses.Response {

	if paymentData.PaymentStatus == "captured" {
		paymentData.PaymentStatus = "paid"
	}

	err := c.OrderRepo.UpdatePaymentStatus(paymentData)

	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't verify payment", err.Error(), nil)
	}

	return responses.ClientReponse(http.StatusOK, "payment veryfied", nil, nil)

}

func (c OrderUseCase) GetInvoiveData(userID, orderID int) responses.Response {

	invoice := models.InvoiceData{
		CompanyName:    "ReadON",
		CompanyAddress: "123 Main Street, City, Country",
		CompanyContact: "+1 234 5678 910",
	}
	var err error
	invoice.Order, err = c.OrderRepo.GetOrder(userID, orderID)
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't process the request", err.Error(), nil)
	}
	orderItems, err := c.OrderRepo.ListOrderItems(orderID)
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't process the request", err.Error(), nil)
	}
	invoice.Date = invoice.Order.CreatedAt.Format("2006-01-02 15:04 ") + "+5:30"
	var orderItemsListing []models.OrderItemsListing
	for i := range orderItems {
		booklisting, err := c.ProductRepo.GetProduct(int(orderItems[i].BookID))
		if err != nil {
			statusCode, err := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't process the request", err.Error(), nil)
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
	invoice.Address, err = c.AddressRepo.GetAddress(invoice.Order.AdressID, uint(userID))
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't process the request", err.Error(), nil)
	}
	invoice.Total = invoice.Order.TotalPrice + invoice.Order.DeleveryCharge
	fmt.Println("invoive data ,", invoice)
	return responses.ClientReponse(http.StatusOK, "invoice generated", nil, invoice)
}

func (c OrderUseCase) GetChartData(pageDetails models.Pagination) responses.Response {
	var chartData models.ChartData

	start, end := getTimeSpan(pageDetails.Filter)
	fmt.Println("start : ,", start, "end : ", end.Add(-time.Millisecond))

	data, err := c.OrderRepo.GetTotalAmountOfSpan(start, end.Add(-time.Millisecond), interval(pageDetails.Filter))
	if err != nil {
		fmt.Println("err", err)
		return responses.ClientReponse(http.StatusInternalServerError, "couldn't process the request", err.Error(), nil)
	}
	fmt.Println(data)
	chartData.Span = []time.Time{start, end.Add(-time.Millisecond)}
	chartData.Data = data
	return responses.ClientReponse(http.StatusOK, "chart data fetched", nil, chartData)
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

func (c OrderUseCase) GetTopTenCategory(filter models.Pagination) responses.Response {
	data, err := c.OrderRepo.FindTopTenCategories(filter)
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't process the request", err.Error(), nil)
	}
	return responses.ClientReponse(http.StatusOK, "top ten categories fetched", nil, data)
}

func (c OrderUseCase) GetTopTenBooks(filter models.Pagination) responses.Response {
	data, err := c.OrderRepo.FindTopTenBooks(filter)
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't process the request", err.Error(), nil)
	}
	var list []models.ListingBook
	copier.Copy(&list, data)
	fmt.Println("data :", data)
	fmt.Println("lsit :", list)
	return responses.ClientReponse(http.StatusOK, "top ten categories fetched", nil, data)
}

func (c OrderUseCase) MakeOrderResponse(order domain.Order) responses.Response {
	var orderListing models.OrdersListing
	copier.Copy(&orderListing, &order)
	orderItems, err := c.OrderRepo.ListOrderItems(int(order.ID))
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch order", err.Error(), nil)
	}
	var orderItemsListing []models.OrderItemsListing
	for i := range orderItems {
		booklisting, err := c.ProductRepo.GetProduct(int(orderItems[i].BookID))
		if err != nil {
			statusCode, err := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't fetch order", err.Error(), nil)
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
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch order", err.Error(), nil)
	}
	if order.PaymentMethodID == 1 {
		orderListing.PaymentMethod = "COD"
	}
	if order.PaymentMethodID == 2 {
		orderListing.PaymentMethod = "Online"
	}
	copier.Copy(&orderListing.Address, address)
	orderListing.Items = orderItemsListing
	return responses.ClientReponse(http.StatusOK, "order fetched", nil, orderListing)
}
