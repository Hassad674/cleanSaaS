package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
	"github.com/hassad/boilerplateSaaS/backend/pkg/hash"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
)

type Service struct {
	users         repository.UserRepository
	resets        repository.PasswordResetRepository
	verifications repository.EmailVerificationRepository
	email         service.EmailService
	jwtMaker      *jwt.Maker
	frontendURL   string
}

func NewService(
	users repository.UserRepository,
	resets repository.PasswordResetRepository,
	verifications repository.EmailVerificationRepository,
	email service.EmailService,
	jwtMaker *jwt.Maker,
	frontendURL string,
) *Service {
	return &Service{
		users:         users,
		resets:        resets,
		verifications: verifications,
		email:         email,
		jwtMaker:      jwtMaker,
		frontendURL:   frontendURL,
	}
}

func (s *Service) Register(ctx context.Context, email, name, password string) (*user.User, string, error) {
	_, err := s.users.FindByEmail(ctx, email)
	if err == nil {
		return nil, "", domain.ErrAlreadyExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, "", err
	}

	hashed, err := hash.Password(password)
	if err != nil {
		return nil, "", err
	}

	u, err := user.New(email, name, hashed)
	if err != nil {
		return nil, "", err
	}

	if err := s.users.Create(ctx, u); err != nil {
		return nil, "", err
	}

	// Send verification email asynchronously (don't block registration)
	if s.email != nil && s.verifications != nil {
		_ = s.sendVerificationEmail(ctx, u)
	}

	token, err := s.jwtMaker.Generate(u.ID, string(u.Role))
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*user.User, string, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", domain.ErrUnauthorized
	}

	if !hash.Check(password, u.PasswordHash) {
		return nil, "", domain.ErrUnauthorized
	}

	token, err := s.jwtMaker.Generate(u.ID, string(u.Role))
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}

func (s *Service) OAuthCallback(ctx context.Context, oauthUser *service.OAuthUser) (*user.User, string, error) {
	u, err := s.users.FindByProvider(ctx, oauthUser.Provider, oauthUser.ProviderID)
	if errors.Is(err, domain.ErrNotFound) {
		u = user.NewOAuth(oauthUser.Email, oauthUser.Name, oauthUser.Provider, oauthUser.ProviderID)
		u.AvatarURL = oauthUser.AvatarURL
		if err := s.users.Create(ctx, u); err != nil {
			return nil, "", err
		}
	} else if err != nil {
		return nil, "", err
	}

	token, err := s.jwtMaker.Generate(u.ID, string(u.Role))
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}

// ForgotPassword generates a password reset token and sends a reset email.
// If the email is not found, it returns nil to avoid leaking user existence.
func (s *Service) ForgotPassword(ctx context.Context, email, frontendURL string) error {
	u, err := s.users.FindByEmail(ctx, email)
	if errors.Is(err, domain.ErrNotFound) {
		// Don't leak whether the user exists
		return nil
	}
	if err != nil {
		return fmt.Errorf("finding user by email: %w", err)
	}

	// Generate a crypto-secure random token (32 bytes = 64 hex chars)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("generating reset token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Store the reset token with 1 hour expiry
	pr := &repository.PasswordReset{
		UserID:    u.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
		CreatedAt: time.Now(),
	}
	if err := s.resets.Create(ctx, pr); err != nil {
		return fmt.Errorf("storing password reset: %w", err)
	}

	// Use the provided frontendURL, or fall back to the configured one
	baseURL := frontendURL
	if baseURL == "" {
		baseURL = s.frontendURL
	}
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)

	// Send the reset email
	if s.email != nil {
		if err := s.email.SendTemplate(ctx, u.Email, "password_reset", map[string]string{
			"name": u.Name,
			"link": resetLink,
		}); err != nil {
			return fmt.Errorf("sending password reset email: %w", err)
		}
	}

	return nil
}

// VerifyEmail validates a verification token and marks the user's email as verified.
func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	ev, err := s.verifications.FindByToken(ctx, token)
	if err != nil {
		return domain.ErrInvalidToken
	}

	if time.Now().After(ev.ExpiresAt) {
		return domain.ErrExpiredToken
	}

	u, err := s.users.FindByID(ctx, ev.UserID)
	if err != nil {
		return fmt.Errorf("finding user for email verification: %w", err)
	}

	if u.EmailVerified {
		// Already verified — clean up token and return success
		_ = s.verifications.DeleteByUserID(ctx, u.ID)
		return nil
	}

	u.VerifyEmail()
	if err := s.users.Update(ctx, u); err != nil {
		return fmt.Errorf("updating user email verified: %w", err)
	}

	// Clean up all verification tokens for this user
	if err := s.verifications.DeleteByUserID(ctx, u.ID); err != nil {
		return fmt.Errorf("deleting verification tokens: %w", err)
	}

	return nil
}

// ResendVerification generates a new verification token and sends it to the user.
func (s *Service) ResendVerification(ctx context.Context, userID string) error {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("finding user for resend verification: %w", err)
	}

	if u.EmailVerified {
		return nil // Already verified, no-op
	}

	// Delete old tokens
	if err := s.verifications.DeleteByUserID(ctx, u.ID); err != nil {
		return fmt.Errorf("deleting old verification tokens: %w", err)
	}

	return s.sendVerificationEmail(ctx, u)
}

// sendVerificationEmail generates a token and sends a verification email.
func (s *Service) sendVerificationEmail(ctx context.Context, u *user.User) error {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("generating verification token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	ev := &repository.EmailVerification{
		UserID:    u.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := s.verifications.Create(ctx, ev); err != nil {
		return fmt.Errorf("storing email verification: %w", err)
	}

	verifyLink := fmt.Sprintf("%s/verify-email?token=%s", s.frontendURL, token)

	if err := s.email.SendTemplate(ctx, u.Email, "verification", map[string]string{
		"name": u.Name,
		"link": verifyLink,
	}); err != nil {
		return fmt.Errorf("sending verification email: %w", err)
	}

	return nil
}

// ResetPassword validates a reset token and updates the user's password.
func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	pr, err := s.resets.FindByToken(ctx, token)
	if err != nil {
		return domain.ErrInvalidToken
	}

	if pr.Used {
		return domain.ErrInvalidToken
	}

	if time.Now().After(pr.ExpiresAt) {
		return domain.ErrExpiredToken
	}

	// Hash the new password
	hashed, err := hash.Password(newPassword)
	if err != nil {
		return fmt.Errorf("hashing new password: %w", err)
	}

	// Find the user and update password
	u, err := s.users.FindByID(ctx, pr.UserID)
	if err != nil {
		return fmt.Errorf("finding user for password reset: %w", err)
	}

	u.PasswordHash = hashed
	u.UpdatedAt = time.Now()
	if err := s.users.Update(ctx, u); err != nil {
		return fmt.Errorf("updating user password: %w", err)
	}

	// Mark the token as used
	if err := s.resets.MarkUsed(ctx, pr.ID); err != nil {
		return fmt.Errorf("marking reset token as used: %w", err)
	}

	return nil
}
