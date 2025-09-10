package lib

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type IdempotencyStore struct {
	client *redis.Client
	ttl    time.Duration
}

type IdempotencyResponse struct {
	StatusCode int               `json:"status_code"`
	Body       json.RawMessage   `json:"body"`
	Headers    map[string]string `json:"headers"`
}

func NewIdempotencyStore(client *redis.Client, ttl time.Duration) *IdempotencyStore {
	return &IdempotencyStore{
		client: client,
		ttl:    ttl,
	}
}

// Store stores a response for an idempotency key
func (s *IdempotencyStore) Store(ctx context.Context, key string, userID string, statusCode int, body interface{}, headers map[string]string) error {
	if s.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	redisKey := fmt.Sprintf("idempotency:%s", key)

	// Initialize headers map if nil
	if headers == nil {
		headers = make(map[string]string)
	}

	// Create response object
	response := IdempotencyResponse{
		StatusCode: statusCode,
		Headers:    headers,
	}

	// Handle body - if it's already bytes, use as is, otherwise marshal to JSON
	var bodyBytes []byte
	var err error

	if body == nil {
		bodyBytes = []byte("{}")
	} else if rawBytes, ok := body.([]byte); ok {
		bodyBytes = rawBytes
	} else {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal response body: %w", err)
		}
	}
	response.Body = bodyBytes

	// Marshal entire response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal idempotency response: %w", err)
	}

	// Store in Redis with TTL
	return s.client.Set(ctx, redisKey, responseBytes, s.ttl).Err()
}

// Get retrieves a stored response for an idempotency key
func (s *IdempotencyStore) Get(ctx context.Context, key string) (*IdempotencyResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}

	redisKey := fmt.Sprintf("idempotency:%s", key)

	val, err := s.client.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, nil // Key not found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get idempotency key: %w", err)
	}

	var response IdempotencyResponse
	err = json.Unmarshal([]byte(val), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal idempotency response: %w", err)
	}

	return &response, nil
}

// Delete removes a stored idempotency key
func (s *IdempotencyStore) Delete(ctx context.Context, key string) error {
	if s.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	redisKey := fmt.Sprintf("idempotency:%s", key)
	return s.client.Del(ctx, redisKey).Err()
}

// GenerateKey generates a unique idempotency key based on user ID and request
func GenerateIdempotencyKey(userID string, method string, path string, body []byte) string {
	// Create a SHA256 hash of the request details (more secure than MD5)
	hash := sha256.New()
	hash.Write([]byte(userID))
	hash.Write([]byte(method))
	hash.Write([]byte(path))
	hash.Write(body)

	return hex.EncodeToString(hash.Sum(nil))[:32] // Use first 32 chars for shorter key
}

// IsIdempotentRequest checks if a request method should be idempotent
func IsIdempotentRequest(method string) bool {
	switch method {
	case "POST", "PUT", "PATCH", "DELETE":
		return true
	default:
		return false
	}
}
