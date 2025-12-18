package domain

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title            string             `bson:"title" json:"title"`
	Location         string             `bson:"location" json:"location"`
	Date             time.Time          `bson:"date" json:"date"`
	TotalTickets     int                `bson:"total_tickets" json:"total_tickets"`
	AvailableTickets int                `bson:"available_tickets" json:"available_tickets"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}