package repository

import (
	"fmt"
	domain "readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	"time"

	"gorm.io/gorm"
)

type OrderDatabse struct {
	DB *gorm.DB
}

func NewOrdersRepository(db *gorm.DB) interfaces.OrderRepository {
	return &OrderDatabse{
		DB: db,
	}
}

func (c OrderDatabse) CreateOrder(order domain.Order, cart []domain.Cart) error {
	fmt.Println("in db create order  : \n", order)
	tx := c.DB.Begin()
	result := tx.Save(&order)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	for i := range cart {
		orderItem := domain.OrderItems{
			OrderID:  order.ID,
			BookID:   cart[i].BookID,
			Quantity: cart[i].Quantity,
			Price:    cart[i].Price,
		}

		err := tx.Create(&orderItem).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err := tx.Where("user_id = ?", order.UserID).Delete(&domain.Cart{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c OrderDatabse) CancelOrder(orderID, userID int) error {
	err := c.DB.Model(&domain.Order{}).Where("user_id = ? AND id = ?", userID, orderID).Update("status", "cancelled").Error
	return err
}

func (c OrderDatabse) DeleteOrder(orderID, userID int) error {
	err := c.DB.Model(&domain.Order{}).Where("user_id = ? AND id = ?", userID, orderID).Delete(&domain.Order{}).Error
	return err
}

func (c OrderDatabse) ListOrders(UserId int, pageDetails models.Pagination) ([]domain.Order, error) {
	var list []domain.Order
	query := c.DB.Model(&domain.Order{}).Where("user_id = ? AND status != ?", UserId, "cancelled").Offset(pageDetails.Offset).Limit(pageDetails.Size)
	if pageDetails.Filter == 1 {
		query = query.Where("payment_status != ?", "failed")
	}
	if pageDetails.Filter == 2 {
		query = query.Where("payment_status = ?", "failed")
	}
	err := query.Find(&list).Error
	return list, err
}

func (c OrderDatabse) GetOrder(userId, orderId int) (domain.Order, error) {
	var order domain.Order
	err := c.DB.Model(&domain.Order{}).Where("user_id = ? AND id = ?", userId, orderId).First(&order).Error
	return order, err
}

// fetches all orders within a time frame
func (c OrderDatabse) GetAllOrders(start, end time.Time) ([]domain.Order, error) {
	var list []domain.Order
	err := c.DB.Model(&domain.Order{}).Where("created_at >= ? AND created_at < ?", start, end).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, err
}

func (c OrderDatabse) GetTotalAmountOfSpan(startTime, endTime time.Time, interval string) ([]models.SalesResult, error) {

	fmt.Println("st :: ", startTime)
	fmt.Println("et ;; ", endTime)

	var results []models.SalesResult

	// SQL query to generate hourly intervals and get the sum of sales for each interval
	query := `
	WITH intervals AS (
		SELECT generate_series(
			?,
			?,
			?::interval
		) AS start_time
	)
	SELECT
		start_time,
		COALESCE(SUM(o.total_quantity), 0) AS total_quantity,
		COALESCE(SUM(o.total_price), 0) AS total_sales
	FROM
		intervals
	LEFT JOIN
		orders o
	ON
		o.created_at >= start_time AND o.created_at < start_time + ?::interval
	GROUP BY
		start_time
	ORDER BY
		start_time;
	`

	// Execute the query and scan the results into the 'results' slice
	err := c.DB.Raw(query, startTime, endTime, interval, interval).Scan(&results).Error

	if err != nil {
		fmt.Println("Query error:", err)
		return results, nil
	}

	return results, nil
}

func (c OrderDatabse) UpdatePaymentStatus(paymentData models.PaymentVerificationData) error {
	err := c.DB.Model(&domain.Order{}).Where("razor_pay_order_id = ?",
		paymentData.RazorOrderId).Updates(map[string]interface{}{"payment_id": paymentData.RazorPaymentId,
		"payment_status": paymentData.PaymentStatus}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c OrderDatabse) UpdateRazorOrderId(userId, orderId int, RazorOrderId string) error {
	db := c.DB.Model(&domain.Order{}).Where("user_id = ? AND id = ? AND payment_status = ?", userId, orderId, "failed")
	err := db.Update("razor_pay_order_id", RazorOrderId).Error
	if err != nil {
		return err
	}
	err = db.Update("payment_status", "processing").Error
	if err != nil {
		return err
	}
	return nil
}
func (c OrderDatabse) GetFailedRazorOrderIds(userId int) ([]string, error) {
	arr := []string{}
	err := c.DB.Model(domain.Order{}).Distinct("razor_pay_order_id").Where("user_id = ?", userId).Find(&arr).Error
	if err != nil {
		return nil, err
	}
	return arr, nil
}
func (c OrderDatabse) GetPrice(userId, orderId int) (float64, error) {
	price := 0.0
	err := c.DB.Model(domain.Order{}).Select("total_price").Where("user_id = ? AND id = ?", userId, orderId).Find(&price).Error
	if err != nil {
		return 0.0, err
	}
	return price, nil
}

func (c OrderDatabse) CheckPymentStatus(orderId int) (string, error) {
	var paymentStatus string
	err := c.DB.Model(domain.Order{}).Select("payment_status").Where("id = ?", orderId).Find(&paymentStatus).Error
	if err != nil {
		return "", err
	}
	return paymentStatus, nil
}

func (c OrderDatabse) ListOrderItems(orderID int) ([]domain.OrderItems, error) {
	var listOfitems []domain.OrderItems
	err := c.DB.Where("order_id = ?", orderID).Find(&listOfitems).Error

	if err != nil {
		fmt.Println("err :", err)
		return listOfitems, err
	}
	return listOfitems, nil
}

func (c OrderDatabse) FetchRazorOrderID(orderID int) (string, error) {
	var razorPayOrderID string
	err := c.DB.Model(&domain.Order{}).Select("razor_pay_order_id").Where("id = ?", orderID).First(&razorPayOrderID).Error
	if err != nil {
		return "", err
	}
	return razorPayOrderID, nil
}

func (c OrderDatabse) FindTopTenCategories(filter models.Pagination) ([]models.TopTenCategory, error) {
	var categories []models.TopTenCategory
	err := c.DB.Debug().Model(domain.OrderItems{}).Joins("JOIN books ON order_items.book_id = books.id ").
		Joins("JOIN categories ON categories.id = books.category_id").
		Select("categories.name as category_name ,SUM(order_items.quantity) as total_sales").
		Group("category_name").
		Order("total_sales DESC").
		Limit(10).Find(&categories).Error

	if err != nil {
		return categories, nil
	}
	fmt.Println("categories ", categories)
	//sum (order_items.quantity as total sales ).GruopBy (categories ) order by desc limit 10
	return categories, nil
}

func (c OrderDatabse) FindTopTenBooks(filter models.Pagination) ([]domain.Book, error) {
	var result []domain.Book
	err := c.DB.Debug().Model(domain.OrderItems{}).Joins("JOIN books ON order_items.book_id = books.id ").
		Select("books.* ,SUM(order_items.quantity) AS total_sales").
		Group("books.id").
		Order("total_sales DESC").
		Limit(10).Find(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}
