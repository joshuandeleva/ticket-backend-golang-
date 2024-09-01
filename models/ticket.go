package models

import (
	"context"
	"time"
)

type Ticket struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	EventID   uint      `json:"eventId"`
	Event     Event     `json:"event" gorm:"foreignkey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID    uint      `json:"userID" gorm:"foreignkey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Entered   bool      `json:"entered" default:"false"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TicketRepository interface {
	CreateOne(context context.Context, userID uint, ticket *Ticket) (*Ticket, error)
	GetMany(context context.Context, userID uint) ([]*Ticket, error)
	GetOne(context context.Context, userID uint, ticketId uint) (*Ticket, error)
	UpdateOne(context context.Context, userID uint, ticketId uint, updateData map[string]interface{}) (*Ticket, error)
	DeleteOne(context context.Context, userID uint, ticketId uint) error
}

type ValidateTicket struct {
	TicketId uint `json:"ticketId"`
	OwnerId   uint `json:"ownerId"`
}
