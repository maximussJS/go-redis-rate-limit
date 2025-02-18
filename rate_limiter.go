package go_redis_rate_limit

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

const luaScriptPath = "rate_limiter.lua"

type RateLimiter struct {
	client           *redis.Client
	script           *redis.Script
	maxAllowedHits   string
	hitsTtlInSeconds string
}

func NewRateLimiter(client *redis.Client, maxAllowedHits, hitsTtlInSeconds int) (*RateLimiter, error) {
	data, err := os.ReadFile(luaScriptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lua script file: %w", err)
	}

	script := redis.NewScript(string(data))

	return &RateLimiter{
		client:           client,
		script:           script,
		maxAllowedHits:   strconv.Itoa(maxAllowedHits),
		hitsTtlInSeconds: strconv.Itoa(hitsTtlInSeconds),
	}, nil
}

func (rl *RateLimiter) Allow(ctx context.Context, tracker string) (bool, error) {
	key := getTrackerKey(tracker)

	result, err := rl.script.Run(ctx, rl.client, []string{key}, rl.maxAllowedHits, rl.hitsTtlInSeconds).Result()

	if err != nil {
		if err == redis.Nil {
			return true, nil
		}

		return false, fmt.Errorf("failed to execute Lua script: %w", err)
	}

	v, ok := result.(int64)
	if !ok {
		return false, fmt.Errorf("unexpected result type from Lua script: %T", v)
	}

	return v == 0, nil
}

func getTrackerKey(tracker string) string {
	return fmt.Sprintf("go-rate-limit-redis:tracker:%s", tracker)
}
