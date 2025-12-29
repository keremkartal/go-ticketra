package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	pb "github.com/keremkartal/goticketra/api/proto/booking"
	"github.com/keremkartal/goticketra/cmd/api-gateway/client"
)

type BookingHandler struct {
	client client.BookingServiceClient
}

func NewBookingHandler(client client.BookingServiceClient) *BookingHandler {
	return &BookingHandler{client: client}
}

type CreateBookingRequest struct {
	UserID      string `json:"user_id"`
	EventID     string `json:"event_id"`
	TicketCount int32  `json:"ticket_count"`
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	var req CreateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri formatı"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pbReq := &pb.CreateBookingRequest{
		UserId:      req.UserID,
		EventId:     req.EventID,
		TicketCount: req.TicketCount,
	}

	res, err := h.client.Client.CreateBooking(ctx, pbReq)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Booking servisi hatası: " + err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(res)
}