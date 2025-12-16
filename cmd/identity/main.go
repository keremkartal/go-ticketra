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
	// 1. Config Yükle
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config yüklenemedi: %v", err)
	}

	// 2. Veritabanına Bağlan
	db := database.ConnectToPostgres(cfg)

	// 3. Otomatik Migrasyon (Tabloları oluştur)
	// Not: Prodüksiyonda bu genelde manuel yapılır ama geliştirme için harikadır.
	db.AutoMigrate(&domain.User{})

	// 4. Katmanları Bağla (Dependency Injection)
	// Repository -> Service -> Handler
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)

	// 5. Fiber Uygulamasını Başlat
	app := fiber.New()

	// 6. Rotaları Tanımla
	api := app.Group("/api/auth") // Tüm rotalar /api/auth ile başlar
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	// 7. Sunucuyu Ayağa Kaldır
	log.Printf("Identity Service %s portunda çalışıyor...", cfg.ServerPort)
	log.Fatal(app.Listen(cfg.ServerPort))
}