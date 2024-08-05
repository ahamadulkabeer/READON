package handler

import (
	"net/http"
	"readon/pkg/api/responses"
	services "readon/pkg/usecase/interface"
	"strconv"

	"github.com/gin-gonic/gin"
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
// @Summary Add item to cart
// @Description Adds a book to the user's cart. If the book is already in the cart, increments the quantity by 1.
// @Tags Cart
// @Accept json
// @Produce json
// @Param bookId query int true "Book ID"
// @Success 200 {object} responses.Response{data=models.ListCartItem} "Item added to cart"
// @Failure 400 {object} responses.Response{error=string} "Bad Request"
// @Failure 500 {object} responses.Response{error=string} "Internal Server Error"
// @Router /cart [post]
func (cr CartHandler) AddToCart(c *gin.Context) {
	userID := c.GetInt("userId")
	bookIdStr := c.Query("bookId")
	bookID, err := strconv.Atoi(bookIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while getting book id", err.Error(), nil))
		return
	}

	response := cr.CartUseCase.AddItem(userID, bookID)

	c.JSON(response.StatusCode, response)

}

// GetCart godoc
// @Summary Get user's cart
// @Description Retrieve the items in the user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response{data=models.ListCart} "Cart fetched successfully"
// @Failure 400 {object} responses.Response "Bad Request"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /cart [get]
func (cr CartHandler) GetCart(c *gin.Context) {

	userID := c.GetInt("userId")

	response := cr.CartUseCase.GetCart(userID)
	c.JSON(response.StatusCode, response)
}

// UpdateQuantity godoc
// @Summary Update quantity of an item in the cart
// @Description Updates the quantity of a book in the user's cart. If the new quantity is less than 1, removes the item from the cart.
// @Tags Cart
// @Accept json
// @Produce json
// @Param bookId query int true "Book ID"
// @Param quantity query int true "New Quantity"
// @Success 200 {object} responses.Response{data=models.ListCartItem} "Item quantity updated"
// @Failure 400 {object} responses.Response{error=string} "Bad Request"
// @Failure 404 {object} responses.Response{error=string} "Item not found"
// @Failure 500 {object} responses.Response{error=string} "Internal Server Error"
// @Router /cart [PUT]
func (cr CartHandler) UpdateQuantity(c *gin.Context) {
	userID := c.GetInt("userId")
	bookidstr := c.Query("bookId")
	qtystr := c.Query("quantity")
	bookID, err := strconv.Atoi(bookidstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while getting book id", err.Error(), nil))
		return
	}
	qty, err := strconv.Atoi(qtystr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while getting quantity ", err.Error(), nil))
		return
	}

	response := cr.CartUseCase.UpdateQty(userID, bookID, qty)

	c.JSON(response.StatusCode, response)
}

// DeleteFromCart godoc
// @Summary Delete item from cart
// @Description Remove an item from the user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param bookId query int true "Book ID"
// @Success 200 {object} responses.Response{data=models.ListCart} "Item removed from the cart"
// @Failure 400 {object} responses.Response "Bad Request"
// @Failure 404 {object} responses.Response "Not Found"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /cart [delete]
func (cr CartHandler) DeleteFromCart(c *gin.Context) {
	userID := c.GetInt("userId")
	bookIdStr := c.Query("bookId")
	bookID, err := strconv.Atoi(bookIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while getting book id", err.Error(), nil))
		return
	}

	response := cr.CartUseCase.DeleteItem(userID, bookID)

	c.JSON(response.StatusCode, response)
}
