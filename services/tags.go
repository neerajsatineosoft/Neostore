package services

import (
	"github.com/Neostore/dbconnection"
	"github.com/Neostore/models"
)

func FetchAllTags() ([]models.Tag, error) {
	database := dbconnection.GetDb()
	var tags []models.Tag
	err := database.Preload("Images", "tag_id IS NOT NULL").Find(&tags).Error
	return tags, err
}
