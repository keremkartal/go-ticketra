package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/keremkartal/goticketra/internal/identity/domain"
	"github.com/keremkartal/goticketra/internal/identity/handler"
	"github.com/keremkartal/goticketra/internal/identity/repository"
	"github.com/keremkartal/goticketra/internal/identity/service"
	"github.com/keremkartal/goticketra/pkg/config"
	"github.com/keremkartal/goticketra/pkg/database"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config yüklenemedi: %v", err)
	}

	db := database.ConnectToPostgres(cfg)

	db.AutoMigrate(&domain.User{})

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)

	app := fiber.New()

	api := app.Group("/api/auth")
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	log.Printf("Identity Service %s portunda çalisiyor...", cfg.IdentityPort)
	log.Fatal(app.Listen(cfg.IdentityPort))
}