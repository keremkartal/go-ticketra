package domain

import (
	"time"
)

type Booking struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      string    `gorm:"not null"`
	EventID     string    `gorm:"not null"` 
	TicketCount int       `gorm:"not null"`
	Status      string    `gorm:"default:'PENDING'"` 
	CreatedAt   time.Time
}