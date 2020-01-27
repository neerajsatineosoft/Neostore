package services

import (
	dbconnection "github.com\Neostore\dbconnection"
)

func CreateOne(data interface{}) error {
	database := dbconnection.GetDb()
	err := database.Create(data).Error
	return err
}

func SaveOne(data interface{}) error {
	database := dbconnection.GetDb()
	err := database.Save(data).Error
	return err
}
