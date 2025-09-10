package services

import (
	"context"

	"foodflow/internal/core"
	"foodflow/internal/core/repos"
	"foodflow/internal/lib"
)

type AuthService struct {
	userRepo    repos.UserRepo
	profileRepo repos.ProfileRepo
	jwtManager  *lib.JWTManager
	hasher      *lib.HashManager
}

func NewAuthService(userRepo repos.UserRepo, profileRepo repos.ProfileRepo, jwtManager *lib.JWTManager) core.AuthService {
	return &AuthService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		jwtManager:  jwtManager,
		hasher:      lib.NewHashManager(),
	}
}

func (s *AuthService) Signup(ctx context.Context, req *core.SignupRequest) (*core.UserResponse, error) {
	// Check if user already exists
	_, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, core.NewConflictError(core.ErrMsgUserAlreadyExists)
	}

	// Hash password
	hashedPassword, err := s.hasher.GenerateHash(req.Password)
	if err != nil {
		return nil, core.NewInternalError("Failed to hash password", err)
	}

	// Create user
	user, err := s.userRepo.CreateUser(ctx, req.Email, hashedPassword, req.Role, core.UserStatusActive)
	if err != nil {
		return nil, core.NewInternalError("Failed to create user", err)
	}

	return &core.UserResponse{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Status: user.Status,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *core.LoginRequest) (*core.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, core.NewUnauthorizedError(core.ErrMsgInvalidCredentials)
	}

	// Check if user is active
	if user.Status != core.UserStatusActive {
		return nil, core.NewUnauthorizedError("Account is suspended")
	}

	// Verify password
	valid, err := s.hasher.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, core.NewInternalError("Failed to verify password", err)
	}
	if !valid {
		return nil, core.NewUnauthorizedError(core.ErrMsgInvalidCredentials)
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID.String(), user.Email, user.Role)
	if err != nil {
		return nil, core.NewInternalError("Failed to generate token", err)
	}

	return &core.AuthResponse{
		AccessToken: token,
		ExpiresIn:   int(s.jwtManager.GetExpirationDuration().Seconds()),
		TokenType:   "Bearer",
	}, nil
}

func (s *AuthService) GetMe(ctx context.Context, userID string) (*core.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &core.UserResponse{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Status: user.Status,
	}, nil
}
