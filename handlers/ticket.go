package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joshuandeleva/go-ticket-backend/models"
	qrcode "github.com/skip2/go-qrcode"
)

type TicketHandler struct {
	repository models.TicketRepository
}

func (h *TicketHandler) GetMany(c *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	userId := uint(c.Locals("userId").(float64))
	tickets, err := h.repository.GetMany(context , userId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status": "fail",
			"error":  err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   tickets,
	})
}

func (h *TicketHandler) GetOne(c *fiber.Ctx) error {
	ticketId , _ := strconv.Atoi(c.Params("ticketId"))
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	userId := uint(c.Locals("userId").(float64))
	ticket , err := h.repository.GetOne(context,userId, uint(ticketId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
		})
	}

	var QRCode []byte
	QRCode, err = qrcode.Encode(
		fmt.Sprintf("ticketId:%v,ownerId:%v", ticketId, userId),
		qrcode.Medium,
		256,
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":  &fiber.Map{
			"ticket": ticket,
			"qrcode": QRCode,
		},
	})
}
func (h *TicketHandler) CreateOne(c *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	userId := uint(c.Locals("userId").(float64))

	ticket := &models.Ticket{}
	if err := c.BodyParser(ticket); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}
	ticket , err := h.repository.CreateOne(context, userId,ticket)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status": "success",
		"message": "Ticket created successfully",
		"data":   ticket,
	})
}

func (h *TicketHandler) ValidateOne(c *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	validateBody := &models.ValidateTicket{}
	if err := c.BodyParser(validateBody); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}
	validateData := make(map[string]interface{})
	validateData["entered"] = true
	ticket , err := h.repository.UpdateOne(context ,validateBody.OwnerId, validateBody.TicketId , validateData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"message": "Ticket valid successfully",
		"data":   ticket,
	})
}

func (h *TicketHandler) UpdateOne(c *fiber.Ctx) error {
	ticketId , _ := strconv.Atoi(c.Params("ticketId"))
	updateData := make(map[string]interface{})
	userId := uint(c.Locals("userId").(float64))

	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}
	ticket , err := h.repository.UpdateOne(context, userId, uint(ticketId), updateData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status": "success",
		"message": "Ticket updated successfully",
		"data":   ticket,
	})
}

func (h *TicketHandler) DeleteOne(c *fiber.Ctx) error {
	ticketId , _ := strconv.Atoi(c.Params("ticketId"))
	userId := uint(c.Locals("userId").(float64))

	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	err := h.repository.DeleteOne(context,userId, uint(ticketId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status": "success",
		"message": "Ticket deleted successfully",
		"data":   nil,
	})
}

func NewTicketHandler(router fiber.Router, repository models.TicketRepository) {
	handler := &TicketHandler{
		repository: repository,
	}
	router.Get("/", handler.GetMany)
	router.Get("/:ticketId", handler.GetOne)
	router.Post("/", handler.CreateOne)
	router.Post("/validate", handler.ValidateOne)
	router.Put("/:ticketId", handler.UpdateOne)
	router.Delete("/:ticketId", handler.DeleteOne)
}