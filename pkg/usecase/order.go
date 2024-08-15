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
	UserRepo    interfaces.UserRepository
}

func NewOrderUseCase(orepo interfaces.OrderRepository,
	crepo interfaces.CartRepository,
	arepo interfaces.AddressRepository,
	prepo interfaces.ProductRepository,
	coupRepo interfaces.CouponRepository,
	UserRepo interfaces.UserRepository) services.OrderUseCase {
	return &OrderUseCase{
		OrderRepo:   orepo,
		CartRepo:    crepo,
		AddressRepo: arepo,
		ProductRepo: prepo,
		CouponRepo:  coupRepo,
		UserRepo:    UserRepo,
	}
}

func (c OrderUseCase) CreateOrder(userID, addressID, PaymentMethodID int, couponCodes []string) responses.Response {

	// fetch cart
	cart, err := c.CartRepo.GetItems(userID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couln't place order", err.Error(), nil)
	}
	if len(cart) == 0 {
		return responses.ClientReponse(http.StatusNotFound, "couln't place order", "cart is empty", nil)
	}

	// calculate total qty and price
	totalQTY := 0
	totalPrice := 0.00
	for _, items := range cart {
		totalQTY += items.Quantity
		totalPrice += items.Price * float64(items.Quantity)
	}

	// initilise order object
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

	// calculate discounts
	statusCode, message, err := helpers.CalculateCouponDiscount(c.CouponRepo, couponCodes, &cart, &order)
	if err != nil {
		return responses.ClientReponse(statusCode, message, err.Error(), nil)
	}
	fmt.Println("order :", order)

	// handle payment
	var razorOrderID string
	if PaymentMethodID == 2 {
		if totalPrice > 1000 {
			return responses.ClientReponse(http.StatusUnprocessableEntity, "couln't place order",
				errors.New("OrderExceedsLimit: Order amount exceeds the maximum allowed limit of 1000"), nil)
		}
		razorOrderID, err = MakeRazorpayPayment(order.DiscountedPrice)
		if err != nil {
			return responses.ClientReponse(http.StatusInternalServerError, "couln't place order", errors.New("payment gate way error").Error(), nil)
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

	// mark coupon as redeemed after placing the order
	for _, coupon := range couponCodes {
		fmt.Println("coupn :", coupon)
		c.CouponRepo.MarkCouponAsRedemed(coupon, order.ID)
	}

	// response with nessessary data
	res := c.MakeOrderResponse(order)
	return responses.ClientReponse(http.StatusOK, "order placed , "+"razorpay order id : "+razorOrderID, nil, res.Data)

}

func (c OrderUseCase) RetryOrder(userID, orderID int) responses.Response {

	//checks the payment status of the order
	paymentStatus, err := c.OrderRepo.CheckPymentStatus(orderID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't retry order", err.Error(), nil)
	}
	if paymentStatus != "failed" {
		return responses.ClientReponse(http.StatusBadRequest, "couldn't retry order", "not a failed payment", nil)
	}

	// fetches the razorpay orderId fro retryign payment
	razorOrderID, err := c.OrderRepo.FetchRazorOrderID(orderID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't retry order", err.Error(), nil)
	}

	// response
	return responses.ClientReponse(http.StatusOK, "order placed",
		nil, "razorpay order id : "+razorOrderID)
}

func (c OrderUseCase) CancelOrder(userID, orderID int) responses.Response {
	//add the paid amount to the wallet
	order, err := c.OrderRepo.GetOrder(userID, orderID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "coulnd't cancel order", err.Error(), nil)
	}

	if order.PaymentStatus == "paid" {
		// send it to a queue for validation and appproval
		c.UserRepo.AddToWallet(order.UserID, order.DiscountedPrice)

	}

	err = c.OrderRepo.CancelOrder(orderID, userID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "coulnd't cancel order", err.Error(), nil)
	}

	err = c.CouponRepo.MarkCouponAsNotRedeemed(uint(orderID))
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "order cancelled , failed to refresh used coupons !", err.Error(), nil)
	}
	// response

	return responses.ClientReponse(http.StatusOK, "order cancelled", nil, nil)
}

func (c OrderUseCase) ListOrders(userID int, pageDetails models.Pagination) responses.Response {
	// setting pagination details
	if pageDetails.Page == 0 {
		pageDetails.Page = 1
	}
	pageDetails.Size = 5
	fmt.Println("page details :", pageDetails)

	// fetches orders
	listOfOrders, err := c.OrderRepo.ListOrders(userID, pageDetails)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch orders", err.Error(), nil)
	}
	if len(listOfOrders) == 0 {
		return responses.ClientReponse(http.StatusNotFound, "no orders found : end of list", nil, nil)
	}

	// populate order data into a listing object
	var orderListing []models.ListOrders
	copier.Copy(&orderListing, &listOfOrders)

	for ind := range orderListing {

		// populating order items into a listing object
		orderItems, err := c.OrderRepo.ListOrderItems(int(orderListing[ind].ID))
		if err != nil {
			statusCode, _ := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't fetch orders", err.Error(), nil)
		}
		var listOfOrderItems []models.ListOrderItems
		copier.Copy(&listOfOrderItems, &orderItems)
		for i := range listOfOrderItems {
			book, err := c.ProductRepo.GetProduct(int(listOfOrderItems[i].BookID))
			if err != nil {
				statusCode, _ := errorhandler.HandleDatabaseError(err)
				return responses.ClientReponse(statusCode, "couldn't fetch orders", err.Error(), nil)
			}
			copier.Copy(&listOfOrderItems[i].Book, &book)
			listOfOrderItems[i].Total += listOfOrderItems[i].Price * float64(listOfOrderItems[i].Quantity)
		}
		orderListing[ind].Items = listOfOrderItems

		// fetch address for each order
		address, err := c.AddressRepo.GetAddress(orderListing[ind].AdressID, uint(userID))
		if err != nil {
			statusCode, _ := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't fetch orders", err.Error(), nil)
		}
		copier.Copy(&orderListing[ind].Address, address)

		// just raw coding instead of ftching from data base
		if listOfOrders[ind].PaymentMethodID == 1 {
			orderListing[ind].PaymentMethod = "COD"
		}
		if listOfOrders[ind].PaymentMethodID == 2 {
			orderListing[ind].PaymentMethod = "Online Payment"
		}
		// rounding
		orderListing[ind].TotalPrice = orderListing[ind].TotalPrice.Round(2)
		orderListing[ind].DiscountedPrice = orderListing[ind].DiscountedPrice.Round(2)
		orderListing[ind].TotalDiscount = orderListing[ind].TotalDiscount.Round(2)
		orderListing[ind].DeleveryCharge = orderListing[ind].DeleveryCharge.Round(2)
	}

	// response
	return responses.ClientReponse(http.StatusOK, "orders fetched successfully", nil, models.PaginatedListOrders{
		Orders:     orderListing,
		Pagination: pageDetails,
	})

}

func (c OrderUseCase) GetOrder(userID, orderID int) responses.Response {

	// fetches order data
	order, err := c.OrderRepo.GetOrder(userID, orderID)
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch order", err.Error(), nil)
	}
	// makes appropriate response
	return c.MakeOrderResponse(order)

}

// to do  : make a pagination here
func (c OrderUseCase) GetAllOrders(filter string) responses.Response {
	var list []domain.Order
	var err error
	list, err = c.OrderRepo.GetAllOrders(getTimeSpan(filter))
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch all orders", err.Error(), nil)
	}
	return responses.ClientReponse(http.StatusOK, "fetched all orders", nil, list)
}

// to populate api keys from .env through side effect
var razorpayKey, razorpaySecret string

func LoadRazorpayConfig(key, secret string) error {
	if key == "" || secret == "" {
		return errors.New("the razorpay api key is missing")
	}
	razorpayKey = key
	razorpaySecret = secret
	return nil
}

func MakeRazorpayPayment(Price float64) (string, error) {
	fmt.Println("key ", razorpayKey, "  sec  ", razorpaySecret)

	client := razorpay.NewClient(razorpaySecret, razorpayKey)
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
		fmt.Println("err :::", err)
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

	// initilise invoice date object
	invoice := models.InvoiceData{
		CompanyName:    "ReadON",
		CompanyAddress: "123 Main Street, City, Country",
		CompanyContact: "+1 234 5678 910",
	}

	// fetch the order
	var err error
	order, err := c.OrderRepo.GetOrder(userID, orderID)
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't process the request", err.Error(), nil)
	}
	invoice.Date = order.CreatedAt.Format("02-01-2006")
	copier.Copy(&invoice.Order, &order)

	// fetch order items
	orderItems, err := c.OrderRepo.ListOrderItems(orderID)
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't process the request", err.Error(), nil)
	}

	// populate order item details into the a listing object
	var orderItemsListing []models.ListOrderItems
	for i := range orderItems {
		productDetails, err := c.ProductRepo.GetProduct(int(orderItems[i].BookID))
		if err != nil {
			statusCode, err := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't process the request", err.Error(), nil)
		}
		orderItemsListing = append(orderItemsListing, models.ListOrderItems{
			BookID:   int(orderItems[i].BookID),
			Price:    orderItems[i].Price,
			Quantity: orderItems[i].Quantity,
			Total:    orderItems[i].Price * float64(orderItems[i].Quantity),
			Book:     models.ListBook{Title: productDetails.Title},
		})
	}
	invoice.OrderItems = orderItemsListing

	// fetch adsress
	invoice.Address, err = c.AddressRepo.GetAddress(invoice.Order.AdressID, uint(userID))
	if err != nil {
		statusCode, err := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't process the request", err.Error(), nil)
	}

	// calculate further data
	invoice.Order.TotalDiscount = invoice.Order.TotalDiscount.Round(2)
	invoice.Total = invoice.Total.Add(invoice.Order.DiscountedPrice)
	invoice.Total = invoice.Total.Add(invoice.Order.DeleveryCharge)
	//invoice.Total = decimal.NewFromFloat(invoice.Order.TotalPrice + invoice.Order.DeleveryCharge - invoice.Order.TotalDiscount).Round(2)

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

func getTimeSpan(filter string) (time.Time, time.Time) {

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
	case "day":
		todayStart := today
		todayEnd := todayStart.Add(24 * time.Hour)
		return todayStart, todayEnd
	case "week":
		thisWeekStart := today.AddDate(0, 0, -int(today.Weekday()))
		thisWeekEnd := thisWeekStart.AddDate(0, 0, 7)
		return thisWeekStart, thisWeekEnd
	case "month":
		thisMonthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, location)
		thisMonthEnd := thisMonthStart.AddDate(0, 1, 0)
		return thisMonthStart, thisMonthEnd
	case "year":
		thisYearStart := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, location)
		thisYearEnd := thisYearStart.AddDate(1, 0, 0)
		return thisYearStart, thisYearEnd
	}
	//baically no span
	return today, today
}

func interval(filter string) string {
	switch filter {
	case "day":
		return "1 hour"
	case "week", "month":
		return "1 day"
	case "year":
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
	var list []models.ListBook
	copier.Copy(&list, data)
	fmt.Println("data :", data)
	fmt.Println("lsit :", list)
	return responses.ClientReponse(http.StatusOK, "top ten categories fetched", nil, data)
}

func (c OrderUseCase) MakeOrderResponse(order domain.Order) responses.Response {
	var orderListing models.ListOrders
	copier.Copy(&orderListing, &order)
	orderItems, err := c.OrderRepo.ListOrderItems(int(order.ID))
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch orders", err.Error(), nil)
	}

	var listOfOrderItems []models.ListOrderItems
	copier.Copy(&listOfOrderItems, &orderItems)
	for i := range listOfOrderItems {
		book, err := c.ProductRepo.GetProduct(int(listOfOrderItems[i].BookID))
		if err != nil {
			statusCode, _ := errorhandler.HandleDatabaseError(err)
			return responses.ClientReponse(statusCode, "couldn't fetch orders", err.Error(), nil)
		}
		copier.Copy(&listOfOrderItems[i].Book, &book)
		listOfOrderItems[i].Total += listOfOrderItems[i].Price * float64(listOfOrderItems[i].Quantity)
	}
	orderListing.Items = listOfOrderItems

	// fetch address for each order
	address, err := c.AddressRepo.GetAddress(orderListing.AdressID, order.UserID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch orders", err.Error(), nil)
	}
	copier.Copy(&orderListing.Address, address)

	// just raw coding instead of ftching from data base
	if order.PaymentMethodID == 1 {
		orderListing.PaymentMethod = "COD"
	}
	if order.PaymentMethodID == 2 {
		orderListing.PaymentMethod = "Online Payment"
	}

	return responses.ClientReponse(http.StatusOK, "order fetched", nil, orderListing)
}

func (c OrderUseCase) GetDataForPaymentpage(orderID int) models.PaymentPageData {
	order, err := c.OrderRepo.GetOrder(1, orderID)
	if err != nil {
		fmt.Println("err occured : ", err.Error())
		return models.PaymentPageData{}
	}
	return models.PaymentPageData{
		RazorpayOrderID: order.RazorPayOrderID,
		OrderID:         order.ID,
		UserName:        order.User.Name,
		FinalPrice:      order.TotalPrice,
	}
}
