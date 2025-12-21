package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/keremkartal/goticketra/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config yüklenemedi: %v", err)
	}

	app := fiber.New()

	app.Group("/api/auth", func(c *fiber.Ctx) error {
		targetURL := cfg.IdentityServiceURL + c.OriginalURL()
		
		if err := proxy.Do(c, targetURL); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Identity servisine ulaşılamıyor: " + err.Error(),
			})
		}
		return nil
	})

	app.Group("/api/events", func(c *fiber.Ctx) error {
		targetURL := cfg.EventServiceURL + c.OriginalURL()
		
		if err := proxy.Do(c, targetURL); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Event servisine ulaşılamıyor: " + err.Error(),
			})
		}
		return nil
	})

	log.Printf("🚀 API Gateway %s portunda çalışıyor...", cfg.GatewayPort)
	log.Fatal(app.Listen(cfg.GatewayPort))
}