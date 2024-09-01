package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joshuandeleva/go-ticket-backend/models"
)

type EventHandler struct {
	repository models.EventRepository
}

func (h *EventHandler) GetMany(c *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	events, err := h.repository.GetMany(context)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status": "fail",
			"error":  err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   events,
	})
}

func (h *EventHandler) GetOne(c *fiber.Ctx) error {
	eventId , _ := strconv.Atoi(c.Params("eventId"))
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	event , err := h.repository.GetOne(context, uint(eventId))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   event,
	})

}

func (h *EventHandler) CreateOne(c *fiber.Ctx) error {
	event := &models.Event{}
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))

	defer cancel();

	if err := c.BodyParser(event); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}
	event , err := h.repository.CreateOne(context, event)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status": "success",
		"message": "Event created successfully",
		"data":   event,
	})

}
func (h *EventHandler) UpdateOne(c *fiber.Ctx) error {
	eventId , _ := strconv.Atoi(c.Params("eventId"))
	updateData := make(map[string]interface{})
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}
	event , err := h.repository.UpdateOne(context, uint(eventId), updateData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status": "success",
		"message": "Event updated successfully",
		"data":   event,
	})
}

func (h *EventHandler) DeleteOne(c *fiber.Ctx) error {
	eventId , _ := strconv.Atoi(c.Params("eventId"))
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	err := h.repository.DeleteOne(context, uint(eventId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status": "fail",
			"message":err.Error(),
			"data":   nil,
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func NewEventHandler(router fiber.Router, repository models.EventRepository) {
	handler := &EventHandler{
		repository: repository,
	}

	router.Get("/", handler.GetMany)
	router.Post("/", handler.CreateOne)
	router.Get("/:eventId", handler.GetOne)
	router.Put("/:eventId", handler.UpdateOne)
	router.Delete("/:eventId", handler.DeleteOne)

}
