package middleware

import (
	"strconv"

	"foodflow/config"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/redis/go-redis/v9"
)

type RateLimitMiddleware struct {
	rateLimiter *lib.RateLimiter
	config      config.AppConfig
}

func NewRateLimitMiddleware(redisClient *redis.Client, config config.AppConfig) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		rateLimiter: lib.NewRateLimiter(redisClient),
		config:      config,
	}
}

// RateLimitByIP applies rate limiting based on IP address
func (m *RateLimitMiddleware) RateLimitByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := "ip:" + clientIP

		allowed, err := m.rateLimiter.IsAllowed(c.Request.Context(), key, m.config.RateLimitRPS, m.config.RateLimitBurst)
		if err != nil {
			// If rate limiter fails, allow the request but log the error
			log.Error().Err(err).Msg("Rate limiter error")
			c.Next()
			return
		}

		if !allowed {
			// Get remaining tokens and reset time for headers
			remaining, _ := m.rateLimiter.GetRemainingTokens(c.Request.Context(), key, m.config.RateLimitRPS, m.config.RateLimitBurst)
			resetTime, _ := m.rateLimiter.GetResetTime(c.Request.Context(), key, m.config.RateLimitRPS, m.config.RateLimitBurst)

			c.Header("X-RateLimit-Limit", strconv.Itoa(m.config.RateLimitBurst))
			c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

			lib.NewResponder().RateLimited(c, "Rate limit exceeded")
			c.Abort()
			return
		}

		// Add rate limit headers
		remaining, _ := m.rateLimiter.GetRemainingTokens(c.Request.Context(), key, m.config.RateLimitRPS, m.config.RateLimitBurst)
		resetTime, _ := m.rateLimiter.GetResetTime(c.Request.Context(), key, m.config.RateLimitRPS, m.config.RateLimitBurst)

		c.Header("X-RateLimit-Limit", strconv.Itoa(m.config.RateLimitBurst))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		c.Next()
	}
}

// RateLimitByUser applies rate limiting based on authenticated user
func (m *RateLimitMiddleware) RateLimitByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			// If no user ID, fall back to IP-based rate limiting
			m.RateLimitByIP()(c)
			return
		}

		key := "user:" + userID

		allowed, err := m.rateLimiter.IsAllowed(c.Request.Context(), key, m.config.RateLimitRPS, m.config.RateLimitBurst)
		if err != nil {
			log.Error().Err(err).Msg("Rate limiter error")
			c.Next()
			return
		}

		if !allowed {
			remaining, _ := m.rateLimiter.GetRemainingTokens(c.Request.Context(), key, m.config.RateLimitRPS, m.config.RateLimitBurst)
			resetTime, _ := m.rateLimiter.GetResetTime(c.Request.Context(), key, m.config.RateLimitRPS, m.config.RateLimitBurst)

			c.Header("X-RateLimit-Limit", strconv.Itoa(m.config.RateLimitBurst))
			c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

			lib.NewResponder().RateLimited(c, "Rate limit exceeded")
			c.Abort()
			return
		}

		// Add rate limit headers
		remaining, _ := m.rateLimiter.GetRemainingTokens(c.Request.Context(), key, m.config.RateLimitRPS, m.config.RateLimitBurst)
		resetTime, _ := m.rateLimiter.GetResetTime(c.Request.Context(), key, m.config.RateLimitRPS, m.config.RateLimitBurst)

		c.Header("X-RateLimit-Limit", strconv.Itoa(m.config.RateLimitBurst))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		c.Next()
	}
}
