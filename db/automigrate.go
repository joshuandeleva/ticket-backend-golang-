package db

import (
	"github.com/joshuandeleva/go-ticket-backend/models"
	"gorm.io/gorm"
)

func DBMigrator(db *gorm.DB) error {
	return db.AutoMigrate(&models.Event{} , &models.Ticket{} , &models.User{})
}