package repository

import (
	"github.com/keremkartal/goticketra/internal/booking/domain"
	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(booking *domain.Booking) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) CreateBooking(booking *domain.Booking) error {
	return r.db.Create(booking).Error
}