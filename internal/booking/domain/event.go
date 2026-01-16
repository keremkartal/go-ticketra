package domain

type BookingCreatedEvent struct {
	BookingID   uint   `json:"booking_id"`
	UserID      string `json:"user_id"`
	TicketCount int    `json:"ticket_count"`
	EventID     string `json:"event_id"`
}