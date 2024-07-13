package handler

import (
	"net/http"
	"readon/pkg/api/responses"
	"readon/pkg/domain"
	services "readon/pkg/usecase/interface"
	"strconv"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while binding parameters , pls check the input ", err.Error(), nil))
		return
	}

	userID := c.GetInt("userId")
	address.UserID = uint(userID)

	response := cr.AddressUseCase.AddAddress(address)

	c.JSON(response.StatusCode, response)

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
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while binding parameters , pls check the input ", err.Error(), nil))
		return
	}
	addressIdStr := c.Param("addressId")
	addressID, err := strconv.Atoi(addressIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"param : addressId not found", err.Error(), nil))
		return
	}

	userID := c.GetInt("userId")

	address.UserID = uint(userID)
	address.ID = uint(addressID)
	response := cr.AddressUseCase.EditAddress(address)
	c.JSON(response.StatusCode, response)
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
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"param : addressId not found", err.Error(), nil))
		return
	}
	userID := c.GetInt("userId")

	response := cr.AddressUseCase.DeleteAddress(uint(addressID), uint(userID))
	c.JSON(response.StatusCode, response)
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
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"param : addressId not found", err.Error(), nil))
		return
	}
	userID := c.GetInt("userId")
	response := cr.AddressUseCase.GetAddress(uint(addressID), uint(userID))

	c.JSON(response.StatusCode, response)
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

	response := cr.AddressUseCase.ListAddress(uint(userID))

	c.JSON(response.StatusCode, response)
}
