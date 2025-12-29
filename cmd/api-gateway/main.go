package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/keremkartal/goticketra/cmd/api-gateway/client"  
	"github.com/keremkartal/goticketra/cmd/api-gateway/handler"
	"github.com/keremkartal/goticketra/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config yüklenemedi: %v", err)
	}

	app := fiber.New()

	bookingClient := client.InitBookingServiceClient(cfg)
	bookingHandler := handler.NewBookingHandler(bookingClient)


	app.Group("/api/auth", func(c *fiber.Ctx) error {
		targetURL := cfg.IdentityServiceURL + c.OriginalURL()
		return proxy.Do(c, targetURL)
	})

	app.Group("/api/events", func(c *fiber.Ctx) error {
		targetURL := cfg.EventServiceURL + c.OriginalURL()
		return proxy.Do(c, targetURL)
	})

	app.Post("/api/bookings", bookingHandler.CreateBooking)

	log.Printf("🚀 API Gateway %s portunda çalışıyor...", cfg.GatewayPort)
	log.Fatal(app.Listen(cfg.GatewayPort))
}