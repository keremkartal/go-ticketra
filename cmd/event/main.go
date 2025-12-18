package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/keremkartal/goticketra/internal/event/handler"
	"github.com/keremkartal/goticketra/internal/event/repository"
	"github.com/keremkartal/goticketra/internal/event/service"
	"github.com/keremkartal/goticketra/pkg/config"
	"github.com/keremkartal/goticketra/pkg/database"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config yüklenemedi: %v", err)
	}

	db := database.ConnectToMongoDB(cfg)

	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo)
	eventHandler := handler.NewEventHandler(eventService)

	app := fiber.New()

	api := app.Group("/api/events")
	api.Post("/", eventHandler.Create)
	api.Get("/", eventHandler.GetAll)
	api.Get("/:id", eventHandler.GetOne)

	port := cfg.EventPort
	if port == "" {
		port = ":8081" 
	}

	log.Printf("Event Service %s portunda çalışıyor...", port)
	log.Fatal(app.Listen(port))
}