package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/keremkartal/goticketra/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB(cfg config.Config) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() 
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		cfg.MongoUser,
		cfg.MongoPassword,
		cfg.MongoHost,
		cfg.MongoPort,
	)

	clientOptions := options.Client().ApplyURI(uri)
	
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB bağlantı hatası: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB'ye erişilemiyor: %v", err)
	}

	fmt.Println(" MongoDB bağlantısı başarılı!")
	
	return client.Database(cfg.MongoDBName)
}