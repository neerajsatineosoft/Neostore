package controllers

import (
	"github.com/Neostore/form"
	"github.com/Neostore/middlewares"
	"github.com/Neostore/models"
	"github.com/Neostore/services"
	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"
)

func RegisterAddressesRoutes(router *gin.RouterGroup) {

	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.GET("/addresses", ListAddresses)
		router.POST("/address", CreateAddress)
	}

}

func ListAddresses(c *gin.Context) {

	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 5
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	// userId:= c.Keys["currentUserId"].(uint) // or
	userId := c.MustGet("currentUserId").(uint)
	includeUser := false
	addresses, totalCommentCount := services.FetchAddressesPage(userId, page, pageSize, includeUser)

	c.JSON(http.StatusOK, form.CreateAddressPagedResponse(c.Request, addresses, page, pageSize, totalCommentCount, includeUser))
}

func CreateAddress(c *gin.Context) {

	user := c.MustGet("currentUser").(models.User)

	var json form.CreateAddress
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, form.CreateBadRequestErrorDto(err))
		return
	}
	firstName := json.FirstName
	lastName := json.LastName
	if firstName == "" {
		firstName = user.FirstName
	}
	if lastName == "" {
		lastName = user.LastName
	}
	address := models.Address{
		FirstName:     firstName,
		LastName:      lastName,
		Country:       json.Country,
		City:          json.City,
		StreetAddress: json.StreetAddress,
		PinCode:       json.PinCode,
		User:          user,
		UserId:        user.ID,
	}

	if err := services.SaveOne(&address); err != nil {
		c.JSON(http.StatusUnprocessableEntity, form.CreateDetailedErrorDto("database_error", err))
		return
	}

	c.JSON(http.StatusOK, form.GetAddressCreatedDto(&address, false))
}
