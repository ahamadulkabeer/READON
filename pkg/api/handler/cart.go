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
// @Summary Add productto cart
// @Description Add a  product cart , if already exist qty++
// @Tags cart
// @Produce json
// @Param bookid formData uint true "Product id"
// @Param userid formData uint true "user id"
// @Success 200 {string} string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/addtocart [post]
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

	userID := c.GetInt("userId")

	response := cr.CartUseCase.GetCart(userID)
	c.JSON(response.StatusCode, response)
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
