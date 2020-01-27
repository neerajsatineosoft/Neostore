package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/Neostore/form"
	"github.com/Neostore/services"
	"net/http"
)

func RegisterPageRoutes(router *gin.RouterGroup) {
	router.GET("", Home)
	router.GET("/home", Home)

}

func Home(c *gin.Context) {

	tags, err := services.FetchAllTags()
	categories, err := services.FetchAllCategories()
	if err != nil {
		c.JSON(http.StatusNotFound, form.CreateDetailedErrorDto("comments", errors.New("Something went wrong")))
		return
	}

	c.JSON(http.StatusOK, form.CreateHomeResponse(tags, categories))
}
