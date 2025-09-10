package middleware

import (
	"bytes"
	"io"
	"net/http"

	"foodflow/config"
	"foodflow/internal/core"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type IdempotencyMiddleware struct {
	idempotencyStore *lib.IdempotencyStore
	config           config.AppConfig
}

func NewIdempotencyMiddleware(redisClient *redis.Client, config config.AppConfig) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{
		idempotencyStore: lib.NewIdempotencyStore(redisClient, config.IdempotencyTTL),
		config:           config,
	}
}

// Idempotency ensures that duplicate requests are handled idempotently
func (m *IdempotencyMiddleware) Idempotency() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to mutating methods
		if !lib.IsIdempotentRequest(c.Request.Method) {
			c.Next()
			return
		}

		// Get idempotency key from header
		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			c.Next()
			return
		}

		// Get user ID for key generation
		userID, exists := GetUserID(c)
		if !exists {
			// If no user ID, use IP address
			userID = c.ClientIP()
		}

		// Read request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			lib.NewResponder().Error(c, core.NewInternalError("Failed to read request body", err))
			c.Abort()
			return
		}

		// Restore request body for downstream handlers
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Generate idempotency key
		key := lib.GenerateIdempotencyKey(userID, c.Request.Method, c.Request.URL.Path, body)

		// Check if we have a cached response
		cachedResponse, err := m.idempotencyStore.Get(c.Request.Context(), key)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get idempotency key")
			c.Next()
			return
		}

		if cachedResponse != nil {
			// Return cached response
			for header, value := range cachedResponse.Headers {
				c.Header(header, value)
			}
			c.Data(cachedResponse.StatusCode, "application/json", cachedResponse.Body)
			c.Abort()
			return
		}

		// Capture response for caching
		responseWriter := &responseCapture{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
			headers:        make(map[string]string),
			statusCode:     200, // Default to 200, will be overridden if WriteHeader is called
		}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Ensure we capture final headers if WriteHeader wasn't called explicitly
		if len(responseWriter.headers) == 0 {
			for key, values := range responseWriter.Header() {
				if len(values) > 0 {
					responseWriter.headers[key] = values[0]
				}
			}
		}

		// Cache successful responses (2xx status codes)
		if responseWriter.statusCode >= 200 && responseWriter.statusCode < 300 {
			// Parse the response body as JSON before storing
			var responseBody interface{}
			if len(responseWriter.body.Bytes()) > 0 {
				responseBody = responseWriter.body.Bytes()
			} else {
				responseBody = map[string]interface{}{}
			}
			
			err = m.idempotencyStore.Store(
				c.Request.Context(),
				key,
				userID,
				responseWriter.statusCode,
				responseBody,
				responseWriter.headers,
			)
			if err != nil {
				log.Error().Err(err).Msg("Failed to store idempotency response")
			}
		}
	}
}

// responseCapture captures the response for idempotency caching
type responseCapture struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	headers    map[string]string
	statusCode int
}

func (w *responseCapture) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *responseCapture) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	// Capture headers when status is written
	for key, values := range w.ResponseWriter.Header() {
		if len(values) > 0 {
			w.headers[key] = values[0]
		}
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseCapture) Header() http.Header {
	return w.ResponseWriter.Header()
}
