package handler

import (
	"fmt"
	"net/http"
	"readon/pkg/domain"
	"readon/pkg/models"
	services "readon/pkg/usecase/interface"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type CartHandler struct {
	CartUseCase services.CartUseCase
}

func NewCartHandler(usecase services.CartUseCase) *CartHandler {
	return &CartHandler{
		CartUseCase: usecase,
	}
}

// AddToCart godoc
// @Summary Add productto cart
// @Description Add a  product cart , if already exist qty++
// @Tags cart
// @Produce json
// @Param bookid formData uint true "Product id"
// @Param userid formData uint true "user id"
// @Success 200 {string} string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/addtocart [postUPI Payment Links is not supported in Test M]
func (cr CartHandler) AddToCart(c *gin.Context) {
	var item domain.Cart

	err := c.Bind(&item)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while binding data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	item.Quantity = 1
	fmt.Println("item :", item)
	// userid is to be get from context ...
	err = cr.CartUseCase.AddItem(item, int(item.UserId))
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while adding item  to cart ",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, "item added to the cart")

}

// @Summary Get user's cart
// @Description Get a user's cart by User ID.
// @Tags cart
// @Produce json
// @Param userid query int true "User ID"
// @Success 200 {object} []models.ListCart "List of items in the cart"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/cart [get]
func (cr CartHandler) GetCart(c *gin.Context) {
	var userid int
	useridstr := c.Query("userid")
	userid, err := strconv.Atoi(useridstr)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return

	}

	cart, err := cr.CartUseCase.GetCart(userid)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while getting cart ",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	var carts []models.ListCart
	copier.Copy(&carts, &cart)
	c.JSON(http.StatusOK, carts)
}

// @Summary Update quantity of a product in the cart
// @Description Update the quantity of a product in the cart for a specific user.
// @Tags cart
// @Produce json
// @Param userid query int true "User ID"
// @Param bookid query int true "Product ID"
// @Param quantity query int true "New quantity"
// @Success 200 {string} string "Quantity updated"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/updatecart [put]
func (cr CartHandler) UpdateQuantity(c *gin.Context) {
	useridstr := c.Query("userid")
	bookidstr := c.Query("bookid")
	qtystr := c.Query("quantity")
	userid, err := strconv.Atoi(useridstr)
	bookid, err := strconv.Atoi(bookidstr)
	qty, err := strconv.Atoi(qtystr)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	err = cr.CartUseCase.UpdateQty(userid, bookid, qty)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while getting cart ",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	response := "quantity updated " + qtystr
	c.JSON(http.StatusOK, response)
}

// @Summary Delete a product from the cart
// @Description Delete a product from the cart for a specific user.
// @Tags cart
// @Produce json
// @Param userid query int true "User ID"
// @Param bookid query int true "Product ID"
// @Success 200 {string} string "Item deleted"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/deleteitem [delete]
func (cr CartHandler) DeleteFromCart(c *gin.Context) {
	useridstr := c.Query("userid")
	bookidstr := c.Query("bookid")
	userid, err := strconv.Atoi(useridstr)
	bookid, err := strconv.Atoi(bookidstr)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	err = cr.CartUseCase.DeleteItem(userid, bookid)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "could delete item ",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	response := "item  deleted :" + bookidstr
	c.JSON(http.StatusOK, response)
}
