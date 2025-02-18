package main

import (
	"context"
	"fmt"
	go_redis_rate_limit "github.com/maximussJS/go-rate-limit-redis"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func main() {
	// Initialize the go-redis client.
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	maxAllowedHits, hitsTtlInSeconds := 5, 10 // Only 5 requests per 10 seconds

	limiter, err := go_redis_rate_limit.NewRateLimiter(redisClient, maxAllowedHits, hitsTtlInSeconds)

	if err != nil {
		log.Fatalf("Could not initialize rate limiter: %v", err)
	}

	ctx := context.Background()
	tracker := "192.168.1.0" // The tracker key can be client IP, user ID, etc.

	for i := 1; i <= 20; i++ {
		allowed, err := limiter.Allow(ctx, tracker)
		if err != nil {
			log.Fatalf("Error checking rate limit: %v", err)
		}

		if allowed {
			fmt.Printf("Request #%d allowed\n", i)
		} else {
			fmt.Printf("Request #%d rate limited\n", i)
		}

		time.Sleep(1 * time.Second)
	}
}
