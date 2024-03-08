package usecase

import (
	"errors"
	"fmt"
	"log"
	"readon/pkg/domain"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
	"time"

	"github.com/razorpay/razorpay-go"
)

type OrderUseCase struct {
	OrderRepo interfaces.OrderRepository
	CartRepo  interfaces.CartRepository
}

func NewOrderUseCase(orepo interfaces.OrderRepository, crepo interfaces.CartRepository) services.OrderUseCase {
	return &OrderUseCase{
		OrderRepo: orepo,
		CartRepo:  crepo,
	}
}

func (c OrderUseCase) CreateOrder(userid, addressid, PaymentMethoadid int) (string, error) {
	var cart []domain.Cart
	cart, err := c.CartRepo.GetItems(userid)
	if err != nil {
		return "", err
	}
	if len(cart) < 1 {
		return "", errors.New("cart is empty")
	}
	var discountedPrice float64
	for i := range cart {
		discountedPrice = discountedPrice + cart[i].Price*float64(cart[i].Quantity)
	}
	var order domain.Order

	if PaymentMethoadid == 1 {
		order.PaymentStatus = "unpaid"
		err := c.ConfirmOrder(userid, addressid, PaymentMethoadid, order)
		if err != nil {
			return "", errors.New("couldn't place order")
		}
	}
	if PaymentMethoadid == 2 {
		fmt.Println("works here  use 1")
		razorOrderID, err := MakeRazorpayPayment(discountedPrice)
		if err != nil {
			return "", err
		}
		fmt.Println("razor id  : ", razorOrderID)
		order.PaymentStatus = "processing"
		order.PaymentId = razorOrderID
		err = c.ConfirmOrder(userid, addressid, PaymentMethoadid, order)
		if err != nil {
			return "", errors.New("couldn't place order")
		}
		return razorOrderID, err
	}
	return "", nil
}

func (c OrderUseCase) ConfirmOrder(userid, addressid, PaymentMethoadid int, order domain.Order) error {
	fmt.Println("works here  use2 confirm")
	//var cart []domain.Cart
	cart, err := c.CartRepo.GetItems(userid)
	if err != nil {
		return err
	}
	var discountedPrice float64
	for i := range cart {
		discountedPrice = discountedPrice + cart[i].Price*float64(cart[i].Quantity)
	}
	for i := range cart {
		order.UserId = cart[i].UserId
		order.BookId = cart[i].BookId
		order.Quantity = cart[i].Quantity
		order.PaymentMethoadId = uint(PaymentMethoadid)
		order.AdressId = uint(addressid)
		order.TotalPrice = discountedPrice

		err = c.OrderRepo.CreateOrder(order, userid)
	}
	if err != nil {
		return err
	}
	return nil
}

func (c OrderUseCase) CancelOrder(userid, orderId int) error {
	err := c.OrderRepo.CancelOrder(orderId, userid)
	if err != nil {
		return err
	}
	return nil
}

func (c OrderUseCase) ListOrders(userid int) ([]domain.Order, error) {
	list, err := c.OrderRepo.GetOrders(userid)
	if err != nil {
		return list, err
	}
	return list, err

}

func (c OrderUseCase) GetOrder(userid, orderid int) (domain.Order, error) {
	order, err := c.OrderRepo.GetOrder(userid, orderid)
	if err != nil {
		return order, err
	}
	return order, nil
}

func (c OrderUseCase) GetAllOrders(filter int) ([]domain.Order, error) {
	var list []domain.Order
	var err error
	today := time.Now().UTC()
	todayStart := time.Now().Truncate(24 * time.Hour)
	if filter == 1 {
		todayEnd := todayStart.Add(24 * time.Hour)
		list, err = c.OrderRepo.GetAllOrders(todayStart, todayEnd)
	}
	if filter == 2 {
		thisWeekStart := todayStart.AddDate(0, 0, -int(todayStart.Weekday())+1)
		thisWeekEnd := thisWeekStart.AddDate(0, 0, 7)
		list, err = c.OrderRepo.GetAllOrders(thisWeekStart, thisWeekEnd)
	}
	if filter == 3 {
		thisMonthStart := time.Date(todayStart.Year(), todayStart.Month(), 1, 0, 0, 0, 0, time.UTC)
		thisMonthEnd := thisMonthStart.AddDate(0, 1, 0)
		list, err = c.OrderRepo.GetAllOrders(thisMonthStart, thisMonthEnd)
	}

	if filter == 4 {
		thisYearStart := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		thisYearEnd := thisYearStart.AddDate(1, 0, 0)
		list, err = c.OrderRepo.GetAllOrders(thisYearStart, thisYearEnd)
	}

	if err != nil {
		return nil, err
	}
	fmt.Println("list :", list)
	return list, nil
}

func MakeRazorpayPayment(Price float64) (string, error) {

	//keyID := "rzp_test_j8TA5AjEgPN2dR"
	//keySecret := "Sh8fzwZmuP7iX7FtIYlH9l0V"

	client := razorpay.NewClient("rzp_test_5aKn1qdbpivTnp", "VDupVmjWyAJY1jBY8m3Ewc45")
	amountInPaise := 1 * 100
	orderData := map[string]interface{}{
		"amount":          amountInPaise, // Amount in smallest currency unit (e.g., in paisa for INR)
		"currency":        "INR",
		"receipt":         "order_rcptid_11",
		"notes":           map[string]string{"key": "value"}, // Additional notes if needed
		"payment_capture": 1,                                 // Auto-capture payment after order creation
	}

	// Create an order
	order, err := client.Order.Create(orderData, nil)
	if err != nil {
		log.Fatal("Error creating order:", err)
	}
	orderID := order["id"].(string)
	fmt.Println("Created Order ID:", orderID)
	fmt.Println("Created Order ID:", order["id"])
	return orderID, nil

}

func (c OrderUseCase) VerifyPayment(paymentID string) error {

	err := c.OrderRepo.UpdatePaymentStatus(paymentID)

	if err != nil {
		return err
	}

	return nil

}
