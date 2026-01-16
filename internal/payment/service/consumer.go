package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/keremkartal/goticketra/internal/payment/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentConsumer struct {
	conn *amqp.Connection
}

func NewPaymentConsumer(conn *amqp.Connection) *PaymentConsumer {
	return &PaymentConsumer{conn: conn}
}

func (c *PaymentConsumer) Start() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"payment_queue", 
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, 
		"",     
		true,   
		false, false, false, nil,
	)
	if err != nil {
		return err
	}

	log.Println(" Payment Service kuyruğu dinliyor...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" Yeni Mesaj Alındı: %s", d.Body)
			
			var event domain.BookingCreatedEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf(" JSON Hatası: %v", err)
				continue
			}

			c.processPayment(event)
		}
	}()

	<-forever 
	return nil
}

func (c *PaymentConsumer) processPayment(event domain.BookingCreatedEvent) {
	log.Printf(" Ödeme alınıyor... BookingID: %d, User: %s", event.BookingID, event.UserID)
	
	time.Sleep(2 * time.Second)

	success := true
	if event.TicketCount > 5 {
		success = false
	}

	if success {
		log.Printf(" Ödeme BAŞARILI! Booking %d onaylandı.", event.BookingID)
	} else {
		log.Printf(" Ödeme BAŞARISIZ! (Yetersiz Bakiye vb.)")
		c.compensateBooking(event)
	}
}

func (c *PaymentConsumer) compensateBooking(event domain.BookingCreatedEvent) {
	log.Println(" TELAFİ İŞLEMİ BAŞLATILIYOR...")
	log.Printf(" Booking Servisine İptal İsteği Gönderiliyor -> BookingID: %d", event.BookingID)
	
	log.Println(" Bilet iptal edildi (Simülasyon). Stok geri yüklendi.")
}