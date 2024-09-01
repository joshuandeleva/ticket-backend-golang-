package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	ID                   uint      `json:"id" gorm:"primarykey;autoIncrement"`
	Name                 string    `json:"name"`
	TotalTicketPurchased int64      `json:"totalTicketPurchased" gorm:"-"`
	TotalTicketEntered   int64      `json:"totalTicketEntered" gorm:"-"`
	Location             string    `json:"location"`
	Date                 time.Time `json:"date"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type EventRepository interface {
	GetMany(ctx context.Context) ([]*Event, error)
	GetOne(ctx context.Context, eventId uint) (*Event, error)
	CreateOne(ctx context.Context, event *Event) (*Event, error)
	UpdateOne(ctx context.Context, eventId uint, updateData map[string]interface{}) (*Event, error)
	DeleteOne(ctx context.Context, eventId uint) error
}


func (e *Event) AfterFind(db *gorm.DB) (err error) {
	baseQuery := db.Model(&Ticket{}).Where("event_id = ?", e.ID)
	 if res  := baseQuery.Count(&e.TotalTicketPurchased); res.Error != nil {
		return res.Error
	}
	if res := baseQuery.Where("entered = ?", true).Count(&e.TotalTicketEntered); res.Error != nil {
		return res.Error
	}
	return nil
}