package repository

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/keremkartal/goticketra/internal/booking/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQRepository interface {
	PublishBookingCreated(event domain.BookingCreatedEvent) error
}

type rabbitMQRepository struct {
	conn *amqp.Connection
}

func NewRabbitMQRepository(conn *amqp.Connection) RabbitMQRepository {
	return &rabbitMQRepository{conn: conn}
}

func (r *rabbitMQRepository) PublishBookingCreated(event domain.BookingCreatedEvent) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close() 
	q, err := ch.QueueDeclare(
		"payment_queue", 
		true,            
		false,           
		false,           
		false,           
		nil,             
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx,
		"",     
		q.Name, 
		false,  
		false,  
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			DeliveryMode: amqp.Persistent, 
		})

	if err != nil {
		log.Printf(" Mesaj kuyruğa atılamadı: %v", err)
		return err
	}

	log.Printf(" Mesaj Kuyruğa Gönderildi: %s", q.Name)
	return nil
}