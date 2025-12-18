package repository

import (
	"context"
	"time"

	"github.com/keremkartal/goticketra/internal/event/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository interface {
	Create(event *domain.Event) error
	FindAll(page int, limit int) ([]domain.Event, error)
	FindByID(id string) (*domain.Event, error)
}

type eventRepository struct {
	collection *mongo.Collection
}

func NewEventRepository(db *mongo.Database) EventRepository {
	return &eventRepository{
		collection: db.Collection("events"),
	}
}

func (r *eventRepository) Create(event *domain.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()
	event.AvailableTickets = event.TotalTickets

	res, err := r.collection.InsertOne(ctx, event)
	if err != nil {
		return err
	}

	event.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *eventRepository) FindAll(page int, limit int) ([]domain.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []domain.Event
	if err = cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) FindByID(id string) (*domain.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var event domain.Event
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}