package handler

import (
	"net/http"
	"readon/pkg/api/responses"
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
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"Error while binding parameters , pls check the input ", err.Error()))
		return
	}

	userID := c.GetInt("userId")

	address.UserID = uint(userID)
	err = cr.AddressUseCase.AddAddress(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.RespondWithError(500,
			"Error while creating resource : couldn't add address", err))
		return
	}

	c.JSON(http.StatusCreated, responses.RespondWithSuccess(http.StatusCreated, "address added successfully", nil))

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
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"Error while binding parameters , pls check the input ", err.Error()))
		return
	}
	addressIdStr := c.Param("addressId")
	addressID, err := strconv.Atoi(addressIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"param : addressId not found", err.Error()))
		return
	}

	userID := c.GetInt("userId")

	address.UserID = uint(userID)
	address.ID = uint(addressID)
	err = cr.AddressUseCase.EditAddress(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.RespondWithError(http.StatusInternalServerError,
			"Error while updating  address : couldn't update address", err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK, "address updated successfully", nil))

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
	addressIdStr := c.Param("addressId")
	addressID, err := strconv.Atoi(addressIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"param : addressId not found", err.Error()))
		return
	}
	userID := c.GetInt("userId")

	err = cr.AddressUseCase.DeleteAddress(uint(addressID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.RespondWithError(http.StatusInternalServerError,
			"Error while deleting  address : couldn't delete address", err.Error()))
		return
	}
	c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusNoContent, "address Deleted successfully", nil))
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
	addressIdStr := c.Param("addressId")
	addressID, err := strconv.Atoi(addressIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"param : addressId not found", err.Error()))
		return
	}
	userID := c.GetInt("userId")
	address, err := cr.AddressUseCase.GetAddress(uint(addressID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.RespondWithError(http.StatusInternalServerError,
			"Error while retriving  address ", err.Error()))
		return
	}
	var addressl models.ListAddress
	copier.Copy(&addressl, &address)
	c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK, "address retrived successfully", addressl))
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

	userID := c.GetInt("userId")

	addressess, err := cr.AddressUseCase.ListAddress(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.RespondWithError(http.StatusInternalServerError,
			"Error while retriving  addressess ", err.Error()))
		return
	}
	var addressl []models.ListAddress
	copier.Copy(&addressl, &addressess)
	c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK, "address retrived successfully", addressl))
}
