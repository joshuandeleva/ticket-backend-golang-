package repositories

import (
	"context"
	"fmt"

	"github.com/joshuandeleva/go-ticket-backend/models"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}




func(r *TicketRepository) CreateOne(ctx context.Context, userId uint ,ticket *models.Ticket)(*models.Ticket, error) {
	ticket.UserID = userId
	res := r.db.Model(ticket).Create(ticket)

	if res.Error != nil {
		return nil , res.Error
	}
	return r.GetOne(ctx ,userId, ticket.ID)

}

func (r *TicketRepository) GetMany(context context.Context ,userId uint) ([]*models.Ticket, error) {
	ticket := []*models.Ticket{}
	res := r.db.Model(&models.Ticket{}).Where("user_id=?" , userId).Preload("Event").Order("updated_at desc").Find(&ticket)

	if res.Error!= nil {
		return nil, res.Error
	}

	return ticket, nil
}


func (r *TicketRepository) GetOne(context context.Context,userId uint, ticketId uint) (*models.Ticket, error) {
	ticket := &models.Ticket{}
	fmt.Print(ticketId , userId)
	res := r.db.Model(ticket).Where("id=?",ticketId).Where("user_id = ?" , userId).Preload("Event").First(ticket)
	if res.Error!= nil {
		return nil, res.Error
	}
	return ticket, nil

}

func (r *TicketRepository) UpdateOne(context context.Context, userId uint, ticketId uint, updateData map[string]interface{}) (*models.Ticket, error) {
	ticket := &models.Ticket{}
	res := r.db.Model(ticket).Where("id =?", ticketId).Updates(updateData)
	if res.Error!= nil {
		return nil, res.Error
	}
	return r.GetOne(context,userId, ticketId)

}

func (r *TicketRepository) DeleteOne(context context.Context,userId uint, ticketId uint) error {
	ticket := &models.Ticket{}
	res := r.db.Model(ticket).Where("id =?", ticketId).Where("user_id=?" , userId).Delete(ticket)
	if res.Error!= nil {
		return res.Error
	}
	return nil

}

func NewTicketRepository(db *gorm.DB) models.TicketRepository {
	return &TicketRepository{db:db}
}

