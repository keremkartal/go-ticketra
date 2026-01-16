package main

import (
	"log"
	"net"

	pb "github.com/keremkartal/goticketra/api/proto/booking"
	"github.com/keremkartal/goticketra/internal/booking/domain"
	"github.com/keremkartal/goticketra/internal/booking/handler"
	"github.com/keremkartal/goticketra/internal/booking/repository"
	"github.com/keremkartal/goticketra/internal/booking/service"
	"github.com/keremkartal/goticketra/pkg/config"
	"github.com/keremkartal/goticketra/pkg/database"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config yüklenemedi: %v", err)
	}

	db := database.ConnectToPostgres(cfg)
	rdb := database.ConnectToRedis(cfg)
	rabbitConn := database.ConnectToRabbitMQ(cfg) 
	defer rabbitConn.Close()

	db.AutoMigrate(&domain.Booking{})

	bookingRepo := repository.NewBookingRepository(db)
	redisRepo := repository.NewRedisRepository(rdb)
	rabbitRepo := repository.NewRabbitMQRepository(rabbitConn) 

	bookingService := service.NewBookingService(bookingRepo, redisRepo, rabbitRepo)
	grpcHandler := handler.NewGrpcHandler(bookingService)

	port := cfg.BookingPort
	if port == "" {
		port = ":50051"
	}
	
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Port dinlenemiyor: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, grpcHandler)

	log.Printf(" Booking Service (gRPC) %s portunda çalışıyor...", port)
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC sunucusu başlatılamadı: %v", err)
	}
}