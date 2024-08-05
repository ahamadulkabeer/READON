package handler

import (
	"fmt"
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

// AddAddress godoc
// @Summary Add a new address
// @Description Add a new address to the system. Requires user authentication and a valid address payload.
// @Tags Address
// @Accept x-www-form-urlencoded
// @Produce json
// @Param addressInput formData models.Address true "Address Input"
// @Success 201 {object} responses.Response{data=models.ListAddress} "Created"
// @Failure 400 {object} responses.Response "Bad Request"
// @Failure 401 {object} responses.Response "Unauthorized"
// @Failure 403 {object} responses.Response "Forbidden"
// @Failure 404 {object} responses.Response "Not Found"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /users/addresses [post]
func (cr AddressHandler) AddAddress(c *gin.Context) {
	var addressInput models.Address
	err := c.Bind(&addressInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while binding parameters , pls check the input ", err.Error(), nil))
		return
	}
	var newAddress domain.Address
	copier.Copy(&newAddress, &addressInput)

	fmt.Println("newaddress ,", newAddress)
	userID := c.GetInt("userId")
	newAddress.UserID = uint(userID)

	response := cr.AddressUseCase.AddAddress(newAddress)

	c.JSON(response.StatusCode, response)

}

// UpdateAddress godoc
// @Summary Update an existing address
// @Description Update an address in the system. Requires user authentication and a valid address payload. The address must exist and belong to the authenticated user.
// @Tags Address
// @Accept x-www-form-urlencoded
// @Produce json
// @Param addressId path int true "Address ID"
// @Param addressInput formData models.Address true "Address Input"
// @Success 200 {object} responses.Response{data=models.ListAddress} "Address Updated"
// @Failure 400 {object} responses.Response "Bad Request"
// @Failure 401 {object} responses.Response "Unauthorized"
// @Failure 404 {object} responses.Response "Not Found"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /users/addresses/{addressId} [put]
func (cr AddressHandler) UpdateAddress(c *gin.Context) {
	var addressInput models.Address
	err := c.Bind(&addressInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while binding parameters , pls check the input ", err.Error(), nil))
		return
	}
	var newAddress domain.Address
	copier.Copy(&newAddress, &addressInput)

	fmt.Println("newaddress ,", newAddress)

	addressIdStr := c.Param("addressId")
	addressID, err := strconv.Atoi(addressIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"param : addressId not found", err.Error(), nil))
		return
	}

	userID := c.GetInt("userId")

	newAddress.UserID = uint(userID)
	newAddress.ID = uint(addressID)
	response := cr.AddressUseCase.EditAddress(newAddress)
	c.JSON(response.StatusCode, response)
}

// DeleteAddress godoc
// @Summary Delete an address
// @Description Delete an address from the system. Requires user authentication. The address must exist and belong to the authenticated user.
// @Tags Address
// @Produce json
// @Param addressId path int true "Address ID"
// @Success 200 {object} responses.Response{data=map[string]interface{}} "Address Deleted"
// @Failure 400 {object} responses.Response "Bad Request"
// @Failure 401 {object} responses.Response "Unauthorized"
// @Failure 404 {object} responses.Response "Not Found"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /users/addresses/{addressId} [delete]
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

// GetAddress godoc
// @Summary Get an address
// @Description Retrieve an address from the system. Requires user authentication. The address must exist and belong to the authenticated user.
// @Tags Address
// @Produce json
// @Param addressId path int true "Address ID"
// @Success 200 {object} responses.Response{data=models.ListAddress} "Address Retrieved"
// @Failure 400 {object} responses.Response "Bad Request"
// @Failure 401 {object} responses.Response "Unauthorized"
// @Failure 404 {object} responses.Response "Not Found"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /users/addresses/{addressId} [get]
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

// ListAddress godoc
// @Summary List all addresses for a user.
// @Description Retrieve a list of addresses associated with the authenticated user. Requires user authentication.
// @Tags Address
// @Produce json
// @Success 200 {object} responses.Response{data=[]models.ListAddress} "Addresses Retrieved"
// @Failure 400 {object} responses.Response "Bad Request"
// @Failure 401 {object} responses.Response "Unauthorized"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /users/addresses [get]
func (cr AddressHandler) ListAddress(c *gin.Context) {

	userID := c.GetInt("userId")

	response := cr.AddressUseCase.ListAddress(uint(userID))

	c.JSON(response.StatusCode, response)
}
