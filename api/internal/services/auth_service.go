// Package services provides business logic for the Setokin API.
package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/setokin/api/internal/config"
	"github.com/setokin/api/internal/models"
	"github.com/setokin/api/internal/repositories"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
	"gorm.io/gorm"
)

// RegisterInput holds the data needed to register a new user.
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// LoginInput holds the data needed to log in.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthTokens holds the generated access and refresh tokens.
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// AuthService defines the interface for authentication business logic.
type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*models.User, *AuthTokens, error)
	Login(ctx context.Context, input LoginInput) (*models.User, *AuthTokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*AuthTokens, error)
	Logout(ctx context.Context, refreshToken string) error
	GetCurrentUser(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
	cfg      config.JWTConfig
	secret   string
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo repositories.UserRepository, cfg config.JWTConfig) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
		secret:   cfg.Secret,
	}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*models.User, *AuthTokens, error) {
	// Validate input
	validationErrors := utils.CollectErrors(
		utils.ValidateEmail(input.Email),
		utils.ValidatePassword(input.Password),
		utils.ValidateRequired("full_name", input.FullName),
		utils.ValidateRole(input.Role),
	)
	if len(validationErrors) > 0 {
		return nil, nil, apperrors.ErrValidation.WithDetails(validationErrors)
	}

	// Check duplicate email
	existing, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err == nil && existing != nil {
		return nil, nil, apperrors.ErrDuplicateResource.WithMessage("Email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, nil, apperrors.ErrInternal.WithMessage("Failed to process password")
	}

	// Set default role
	role := input.Role
	if role == "" {
		role = "staff"
	}

	// Create user
	user := &models.User{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		FullName:     input.FullName,
		Role:         role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, apperrors.ErrInternal.WithMessage("Failed to create user")
	}

	// Generate tokens
	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (*models.User, *AuthTokens, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, nil, apperrors.ErrAuthenticationRequired.WithMessage("Invalid email or password")
	}

	// Check if account is active
	if !user.IsActive {
		return nil, nil, apperrors.ErrInsufficientPermissions.WithMessage("Account is inactive")
	}

	// Verify password
	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		return nil, nil, apperrors.ErrAuthenticationRequired.WithMessage("Invalid email or password")
	}

	// Generate tokens
	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	// Validate refresh token JWT
	claims, err := utils.ValidateRefreshToken(refreshToken, s.secret)
	if err != nil {
		return nil, apperrors.ErrTokenInvalid.WithMessage("Invalid or expired refresh token")
	}

	// Check token in database
	tokenHash := hashToken(refreshToken)
	storedToken, err := s.userRepo.FindRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, apperrors.ErrTokenInvalid.WithMessage("Invalid or expired refresh token")
	}

	// Revoke old refresh token
	if err := s.userRepo.RevokeRefreshToken(ctx, storedToken.ID); err != nil {
		return nil, apperrors.ErrInternal
	}

	// Get user
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, apperrors.ErrTokenInvalid
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrTokenInvalid
	}

	// Generate new token pair
	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	storedToken, err := s.userRepo.FindRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		// Token not found or already revoked — that's OK for logout
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return s.userRepo.RevokeRefreshToken(ctx, storedToken.ID)
}

func (s *authService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrNotFound.WithMessage("User not found")
	}
	return user, nil
}

// generateTokenPair creates access + refresh tokens and stores the refresh token hash.
func (s *authService) generateTokenPair(ctx context.Context, user *models.User) (*AuthTokens, error) {
	accessToken, err := utils.GenerateAccessToken(
		user.ID, user.Email, user.Role, s.secret, s.cfg.AccessExpiry,
	)
	if err != nil {
		return nil, apperrors.ErrInternal.WithMessage(fmt.Sprintf("Failed to generate access token: %v", err))
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.secret, s.cfg.RefreshExpiry)
	if err != nil {
		return nil, apperrors.ErrInternal.WithMessage(fmt.Sprintf("Failed to generate refresh token: %v", err))
	}

	// Store refresh token hash
	tokenHash := hashToken(refreshToken)
	storedToken := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: utils.MustParseRefreshExpiry(refreshToken, s.secret),
	}
	if err := s.userRepo.SaveRefreshToken(ctx, storedToken); err != nil {
		return nil, apperrors.ErrInternal.WithMessage("Failed to save refresh token")
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.cfg.AccessExpiry.Seconds()),
	}, nil
}

// hashToken creates a SHA-256 hash of a token string.
func hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
