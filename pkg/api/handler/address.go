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

type AddressHandler struct {
	AddressUseCase services.AddressUsecase
}

func NewAddressHandler(usecase services.AddressUsecase) *AddressHandler {
	return &AddressHandler{
		AddressUseCase: usecase,
	}
}

// @Summary Add address
// @Description Add a new address for a user.
// @Tags address
// @Accept json
// @Produce json
// @Param address body domain.Address true "Address details"
// @Success 200 {string} string "Address added"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/addaddress [post]
func (cr AddressHandler) AddAddress(c *gin.Context) {
	var address domain.Address
	err := c.Bind(&address)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while binding data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	fmt.Println("adress in hdler :", address)

	err = cr.AddressUseCase.AddAddress(address, int(address.UserId))
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while adding Adreess ",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	c.JSON(http.StatusOK, "Adress Added")

}

// @Summary Update address
// @Description Update an existing address for a user.
// @Tags address
// @Accept json
// @Produce json
// @Param address body domain.Address true "Updated address details"
// @Success 200 {string} string "Address updated"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/updateaddress [put]
func (cr AddressHandler) UpdateAddress(c *gin.Context) {
	var address domain.Address
	err := c.Bind(&address)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while binding data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	err = cr.AddressUseCase.EditAddress(address, int(address.UserId))
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while updating Adreess ",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	c.JSON(http.StatusOK, "Adress Updated")

}

// @Summary Delete address
// @Description Delete an address by its ID.
// @Tags address
// @Produce json
// @Param addressid query int true "Address ID"
// @Success 200 {string} string "Address deleted"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/deleteaddress [delete]
func (cr AddressHandler) DeleteAddress(c *gin.Context) {
	addressIdStr := c.Query("addressid")
	addressId, err := strconv.Atoi(addressIdStr)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	err = cr.AddressUseCase.DeleteAddress(addressId)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while deleting Adreess ",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, "Adress Deleted")
}

// @Summary Get address
// @Description Retrieve an address by its ID.
// @Tags address
// @Produce json
// @Param addressid query int true "Address ID"
// @Success 200 {object} models.ListAddress "Address details"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/getaddress [get]
func (cr AddressHandler) GetAddress(c *gin.Context) {
	addressIdStr := c.Query("addressid")
	addressId, err := strconv.Atoi(addressIdStr)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	address, err := cr.AddressUseCase.GetAddress(addressId)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while getting Adreess ",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	var addressl models.ListAddress
	copier.Copy(&addressl, &address)
	c.JSON(http.StatusOK, addressl)
}

// @Summary List addresses
// @Description Retrieve a list of addresses for a user by user ID.
// @Tags address
// @Produce json
// @Param userid query int true "User ID"
// @Success 200 {array} []models.ListAddress "List of addresses"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/listaddresses [get]
func (cr AddressHandler) ListAddress(c *gin.Context) {
	userIdStr := c.Query("userid")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	addressess, err := cr.AddressUseCase.ListAddress(userId)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while getting Adreessess ",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	var addressl []models.ListAddress
	copier.Copy(&addressl, &addressess)
	c.JSON(http.StatusOK, addressl)
}
