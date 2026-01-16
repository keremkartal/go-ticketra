package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/keremkartal/goticketra/internal/booking/domain"
	"github.com/keremkartal/goticketra/internal/booking/repository"
)

type BookingService interface {
	CreateBooking(userID, eventID string, ticketCount int) (*domain.Booking, error)
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	redisRepo   repository.RedisRepository
	rabbitRepo  repository.RabbitMQRepository
}

func NewBookingService(bRepo repository.BookingRepository, rRepo repository.RedisRepository, rabbitRepo repository.RabbitMQRepository) BookingService {
	return &bookingService{
		bookingRepo: bRepo,
		redisRepo:   rRepo,
		rabbitRepo:  rabbitRepo,
	}
}

func (s *bookingService) CreateBooking(userID, eventID string, ticketCount int) (*domain.Booking, error) {
	lockKey := fmt.Sprintf("lock:event:%s", eventID)
	acquired, err := s.redisRepo.AcquireLock(lockKey, 5*time.Second)
	if err != nil { return nil, err }
	if !acquired { return nil, errors.New("şu an bu etkinlik için çok fazla işlem var") }
	defer s.redisRepo.ReleaseLock(lockKey)

	booking := &domain.Booking{
		UserID:      userID,
		EventID:     eventID,
		TicketCount: ticketCount,
		Status:      "PENDING",
		CreatedAt:   time.Now(),
	}

	if err := s.bookingRepo.CreateBooking(booking); err != nil {
		return nil, err
	}

	event := domain.BookingCreatedEvent{
		BookingID:   booking.ID,
		UserID:      booking.UserID,
		TicketCount: booking.TicketCount,
		EventID:     booking.EventID,
	}

	if err := s.rabbitRepo.PublishBookingCreated(event); err != nil {
		fmt.Printf(" DİKKAT: Booking %d için ödeme mesajı gönderilemedi!\n", booking.ID)
	}

	return booking, nil
}