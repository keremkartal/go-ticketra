package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/keremkartal/goticketra/internal/event/service"
)

type EventHandler struct {
	service service.EventService
}

func NewEventHandler(service service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

type CreateEventRequest struct {
	Title        string `json:"title"`
	Location     string `json:"location"`
	Date         string `json:"date"` 
	TotalTickets int    `json:"total_tickets"`
}

func (h *EventHandler) Create(c *fiber.Ctx) error {
	var req CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	event, err := h.service.CreateEvent(req.Title, req.Location, req.Date, req.TotalTickets)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(event)
}

func (h *EventHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	events, err := h.service.ListEvents(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(events)
}

func (h *EventHandler) GetOne(c *fiber.Ctx) error {
	id := c.Params("id")
	event, err := h.service.GetEvent(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Etkinlik bulunamadı"})
	}
	return c.Status(200).JSON(event)
}