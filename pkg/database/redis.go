package database

import (
	"context"
	"fmt"
	"log"

	"github.com/keremkartal/goticketra/pkg/config"
	"github.com/redis/go-redis/v9"
)

func ConnectToRedis(cfg config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       0,             
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis bağlantı hatası: %v", err)
	}

	fmt.Println(" Redis bağlantısı başarılı!")
	return rdb
}