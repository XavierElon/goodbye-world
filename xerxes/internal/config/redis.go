package config

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() {
	log.Println("Initializing Redis connection...")
	
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis123",
		DB:       0,
		// Add connection pool settings
		PoolSize:     10,
		MinIdleConns: 5,
		// Add timeout settings
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test the connection with retry logic
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		
		_, err := RedisClient.Ping(ctx).Result()
		cancel()
		
		if err == nil {
			log.Println("Successfully connected to Redis")
			return
		}
		
		log.Printf("Redis connection attempt %d/%d failed: %v", i+1, maxRetries, err)
		
		if i < maxRetries-1 {
			log.Printf("Retrying in 2 seconds...")
			time.Sleep(2 * time.Second)
		}
	}
	
	log.Fatal("Failed to connect to Redis after all retry attempts")
}