package client

import (
	"log"

	pb "github.com/keremkartal/goticketra/api/proto/booking"
	"github.com/keremkartal/goticketra/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BookingServiceClient struct {
	Client pb.BookingServiceClient
}

func InitBookingServiceClient(cfg config.Config) BookingServiceClient {
	cc, err := grpc.Dial(cfg.BookingPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Booking servisine bağlanılamadı: %v", err)
	}

	log.Println(" Gateway -> Booking Service (gRPC) bağlantısı kuruldu.")

	return BookingServiceClient{
		Client: pb.NewBookingServiceClient(cc),
	}
}