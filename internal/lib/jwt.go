package lib

import (
	"fmt"
	"time"

	"foodflow/config"
	"foodflow/internal/core"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	secret     []byte
	issuer     string
	expiration time.Duration
}

type Claims struct {
	UserID string        `json:"user_id"`
	Email  string        `json:"email"`
	Role   core.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTManager(cfg config.JWTConfig) (*JWTManager, error) {
	if cfg.Secret == "" {
		return nil, fmt.Errorf("JWT secret is required")
	}

	return &JWTManager{
		secret:     []byte(cfg.Secret),
		issuer:     cfg.Issuer,
		expiration: cfg.Expiration,
	}, nil
}

func (j *JWTManager) GenerateToken(userID, email string, role core.UserRole) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID,
			Audience:  []string{"foodflow"},
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expiration)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (j *JWTManager) GetExpirationDuration() time.Duration {
	return j.expiration
}
