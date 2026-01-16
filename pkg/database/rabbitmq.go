package database

import (
	"fmt"
	"log"

	"github.com/keremkartal/goticketra/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMQ(cfg config.Config) *amqp.Connection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitUser,
		cfg.RabbitPass,
		cfg.RabbitHost,
		cfg.RabbitPort,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("RabbitMQ bağlantı hatası: %v", err)
	}

	fmt.Println("RabbitMQ bağlantısı başarılı!")
	return conn
}