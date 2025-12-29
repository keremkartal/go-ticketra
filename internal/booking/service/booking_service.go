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
}

func NewBookingService(bRepo repository.BookingRepository, rRepo repository.RedisRepository) BookingService {
	return &bookingService{
		bookingRepo: bRepo,
		redisRepo:   rRepo,
	}
}

func (s *bookingService) CreateBooking(userID, eventID string, ticketCount int) (*domain.Booking, error) {
	lockKey := fmt.Sprintf("lock:event:%s", eventID)
	
	acquired, err := s.redisRepo.AcquireLock(lockKey, 5*time.Second)
	if err != nil {
		return nil, err
	}
	if !acquired {
		return nil, errors.New("şu an bu etkinlik için çok fazla işlem var, lütfen tekrar deneyin")
	}

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

	return booking, nil
}