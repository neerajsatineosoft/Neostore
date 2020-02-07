package controllers

import (
	"errors"

	dtos "github.com/Neostore/form"
	"github.com/Neostore/services"
	"github.com/gin-gonic/gin"

	"github.com/Neostore/models"

	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUserRoutes(router *gin.RouterGroup) {
	router.POST("/register", UsersRegistration)
	router.POST("/login", UsersLogin)
	//router.POST("/createuser",Createuser)
}

func UsersRegistration(c *gin.Context) {

	var json dtos.RegisterRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	if err := services.CreateOne(&models.User{
		//Username:  json.Username,
		Password:  string(password),
		FirstName: json.FirstName,
		LastName:  json.LastName,
		Email:     json.Email,
		Phoneno:   json.Phoneno,
		Gender:    json.Gender,
	}); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("database", err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success":       true,
		"full_messages": []string{"User created successfully"}})
}

func UsersLogin(c *gin.Context) {

	var json dtos.LoginRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	user, err := services.FindOneUser(&models.User{Email: json.Email})

	if err != nil {
		c.JSON(http.StatusForbidden, dtos.CreateDetailedErrorDto("login_error", err))
		return
	}

	if user.IsValidPassword(json.Password) != nil {
		c.JSON(http.StatusForbidden, dtos.CreateDetailedErrorDto("login", errors.New("invalid credentials")))
		return
	}

	c.JSON(http.StatusOK, dtos.CreateLoginSuccessful(&user))

}
