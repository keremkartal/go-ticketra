package database

import (
	"fmt"
	"log"

	"github.com/keremkartal/goticketra/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectToPostgres, verilen config ile veritabanına bağlanır ve *gorm.DB döner.
func ConnectToPostgres(cfg config.Config) *gorm.DB {
	// DSN (Data Source Name): Bağlantı cümlesi
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	fmt.Println("✅ PostgreSQL bağlantısı başarılı!")
	return db
}