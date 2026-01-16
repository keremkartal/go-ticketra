package main

import (
	"log"

	"github.com/keremkartal/goticketra/internal/payment/service"
	"github.com/keremkartal/goticketra/pkg/config"
	"github.com/keremkartal/goticketra/pkg/database"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config yüklenemedi: %v", err)
	}

	rabbitConn := database.ConnectToRabbitMQ(cfg)
	defer rabbitConn.Close()

	consumer := service.NewPaymentConsumer(rabbitConn)

	log.Println(" Payment Service (Worker) Başlatılıyor...")
	
	if err := consumer.Start(); err != nil {
		log.Fatalf("Consumer hatası: %v", err)
	}
}