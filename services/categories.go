package services

import (
	"github.com/Neostore/dbconnection"
	"github.com/Neostore/models"
)

func FetchAllCategories() ([]models.Category, error) {
	database := dbconnection.GetDb()
	var categories []models.Category
	err := database.Preload("Images", "category_id IS NOT NULL").Find(&categories).Error
	return categories, err
}
