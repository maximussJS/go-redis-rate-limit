# go-rate-limit-redis

A simple, Lua-based rate limiter for Redis written in Go. This package allows you to easily limit the number of requests for a given key (e.g. IP address, user ID, etc.) using Redis for storage and atomicity. The rate limiter uses a Lua script to increment request counters and set expirations, ensuring high performance and accurate rate limiting even under concurrent access.

## Features

- **Atomic Rate Limiting:** Uses a Redis Lua script to ensure atomic increment and expiration.
- **Easy Integration:** Works with [go-redis](https://github.com/redis/go-redis) v9.
- **Customizable Limits:** Specify maximum allowed requests and time window.

## Installation

To install the package, run:

```bash
go get github.com/maximussJS/go-rate-limit-redis
```

## Usage
Below is an example of how to use the rate limiter in your application:


```go
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

	// Set rate limiting parameters:
	// Only allow 5 requests per 10 seconds.
	maxAllowedHits, hitsTtlInSeconds := 5, 10

	// Initialize the rate limiter with the Redis client and rate limiting parameters.
	limiter, err := go_redis_rate_limit.NewRateLimiter(redisClient, maxAllowedHits, hitsTtlInSeconds)
	if err != nil {
		log.Fatalf("Could not initialize rate limiter: %v", err)
	}

	ctx := context.Background()
	// The tracker can be a client IP, user ID, etc.
	tracker := "192.168.1.0"

	// Simulate 20 requests in next 20 seconds (1 request per second).
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
```

## API

### `NewRateLimiter`

```go
func NewRateLimiter(client *redis.Client, maxAllowedHits, hitsTtlInSeconds int) (*RateLimiter, error)
```

Parameters:
- **client**: A pointer to a go-redis v9 redis.Client instance.
- **maxAllowedHits**: The maximum number of allowed requests (hits) within the defined time window.
- **hitsTtlInSeconds**: The duration of the time window in seconds.

Returns:
- A new RateLimiter instance.
- An error if the rate limiter could not be initialized.

### `RateLimiter.Allow`
```go
func (rl *RateLimiter) Allow(ctx context.Context, tracker string) (bool, error)
```

Parameters:
- **ctx**: The context for managing deadlines and cancellation.
- **tracker**: A unique identifier (such as an IP address, user ID, etc.) used to track requests for rate limiting.

Returns:
- A boolean indicating whether the request is allowed. `true` if the request is allowed, `false` if the request is rate limited.
- An error if the rate limit check failed.

# License

MIT License
