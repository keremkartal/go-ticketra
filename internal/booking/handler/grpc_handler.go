package handler

import (
	"context"
	"fmt"

	pb "github.com/keremkartal/goticketra/api/proto/booking"
	"github.com/keremkartal/goticketra/internal/booking/service"
)

type GrpcHandler struct {
	pb.UnimplementedBookingServiceServer 
	service                              service.BookingService
}

func NewGrpcHandler(service service.BookingService) *GrpcHandler {
	return &GrpcHandler{service: service}
}

func (h *GrpcHandler) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	fmt.Printf(" gRPC İsteği Geldi: User=%s, Event=%s\n", req.UserId, req.EventId)

	booking, err := h.service.CreateBooking(req.UserId, req.EventId, int(req.TicketCount))
	
	if err != nil {
		return &pb.CreateBookingResponse{
			Status:       "FAILED",
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.CreateBookingResponse{
		BookingId: fmt.Sprintf("%d", booking.ID), 
		Status:    "SUCCESS",
	}, nil
}