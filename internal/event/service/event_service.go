package service

import (
	"time" 
	"github.com/keremkartal/goticketra/internal/event/domain"
	"github.com/keremkartal/goticketra/internal/event/repository"
)

type EventService interface {
	CreateEvent(title, location string, date string, totalTickets int) (*domain.Event, error)
	ListEvents(page, limit int) ([]domain.Event, error)
	GetEvent(id string) (*domain.Event, error)
}

type eventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{repo: repo}
}

func (s *eventService) CreateEvent(title, location string, dateStr string, totalTickets int) (*domain.Event, error) {
	layout := time.RFC3339
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, err
	}

	event := &domain.Event{
		Title:        title,
		Location:     location,
		Date:         parsedDate,
		TotalTickets: totalTickets,
	}

	err = s.repo.Create(event)
	return event, err
}

func (s *eventService) ListEvents(page, limit int) ([]domain.Event, error) {
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }
	return s.repo.FindAll(page, limit)
}

func (s *eventService) GetEvent(id string) (*domain.Event, error) {
	return s.repo.FindByID(id)
}