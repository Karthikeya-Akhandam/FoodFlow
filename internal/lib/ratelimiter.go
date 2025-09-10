package lib

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
}

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{
		client: client,
	}
}

// TokenBucket implements a token bucket rate limiter
type TokenBucket struct {
	Key        string
	Capacity   int
	Tokens     int
	RefillRate int // tokens per second
	LastRefill time.Time
}

// IsAllowed checks if a request is allowed based on token bucket algorithm
func (rl *RateLimiter) IsAllowed(ctx context.Context, key string, limit, burst int) (bool, error) {
	now := time.Now()

	// Get current bucket state
	bucket, err := rl.getBucket(ctx, key, limit, burst)
	if err != nil {
		return false, err
	}

	// Refill tokens based on time elapsed
	elapsed := now.Sub(bucket.LastRefill)
	tokensToAdd := int(elapsed.Seconds() * float64(bucket.RefillRate))

	if tokensToAdd > 0 {
		bucket.Tokens = min(bucket.Capacity, bucket.Tokens+tokensToAdd)
		bucket.LastRefill = now
	}

	// Check if request is allowed
	if bucket.Tokens > 0 {
		bucket.Tokens--
		err = rl.saveBucket(ctx, key, bucket)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// Save updated bucket state
	err = rl.saveBucket(ctx, key, bucket)
	if err != nil {
		return false, err
	}

	return false, nil
}

// GetRemainingTokens returns the number of remaining tokens
func (rl *RateLimiter) GetRemainingTokens(ctx context.Context, key string, limit, burst int) (int, error) {
	bucket, err := rl.getBucket(ctx, key, limit, burst)
	if err != nil {
		return 0, err
	}

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(bucket.LastRefill)
	tokensToAdd := int(elapsed.Seconds() * float64(bucket.RefillRate))

	if tokensToAdd > 0 {
		bucket.Tokens = min(bucket.Capacity, bucket.Tokens+tokensToAdd)
		bucket.LastRefill = now
	}

	return bucket.Tokens, nil
}

// GetResetTime returns when the bucket will be fully refilled
func (rl *RateLimiter) GetResetTime(ctx context.Context, key string, limit, burst int) (time.Time, error) {
	bucket, err := rl.getBucket(ctx, key, limit, burst)
	if err != nil {
		return time.Time{}, err
	}

	if bucket.Tokens >= bucket.Capacity {
		return time.Now(), nil
	}

	tokensNeeded := bucket.Capacity - bucket.Tokens
	secondsToReset := float64(tokensNeeded) / float64(bucket.RefillRate)

	return bucket.LastRefill.Add(time.Duration(secondsToReset) * time.Second), nil
}

func (rl *RateLimiter) getBucket(ctx context.Context, key string, limit, burst int) (*TokenBucket, error) {
	redisKey := fmt.Sprintf("rate_limit:%s", key)

	// Try to get existing bucket
	val, err := rl.client.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		// Create new bucket
		return &TokenBucket{
			Key:        redisKey,
			Capacity:   burst,
			Tokens:     burst,
			RefillRate: limit,
			LastRefill: time.Now(),
		}, nil
	} else if err != nil {
		return nil, err
	}

	// Parse existing bucket
	bucket := &TokenBucket{}
	var lastRefillUnix int64
	_, err = fmt.Sscanf(val, "%d:%d:%d", &bucket.Tokens, &bucket.RefillRate, &lastRefillUnix)
	if err != nil {
		// If parsing fails, create new bucket
		return &TokenBucket{
			Key:        redisKey,
			Capacity:   burst,
			Tokens:     burst,
			RefillRate: limit,
			LastRefill: time.Now(),
		}, nil
	}

	bucket.Key = redisKey
	bucket.Capacity = burst
	bucket.LastRefill = time.Unix(lastRefillUnix, 0)

	return bucket, nil
}

func (rl *RateLimiter) saveBucket(ctx context.Context, key string, bucket *TokenBucket) error {
	redisKey := fmt.Sprintf("rate_limit:%s", key)
	val := fmt.Sprintf("%d:%d:%d", bucket.Tokens, bucket.RefillRate, bucket.LastRefill.Unix())

	// Set with expiration (24 hours)
	return rl.client.Set(ctx, redisKey, val, 24*time.Hour).Err()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
