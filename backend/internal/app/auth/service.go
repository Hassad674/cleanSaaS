package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/org"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
	"github.com/hassad/boilerplateSaaS/backend/pkg/hash"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
)

type Service struct {
	users           repository.UserRepository
	orgs            repository.OrganizationRepository
	resets          repository.PasswordResetRepository
	verifications   repository.EmailVerificationRepository
	refreshTokens   repository.RefreshTokenRepository
	tx              repository.TxManager
	email           service.EmailService
	jwtMaker        *jwt.Maker
	frontendURL     string
	refreshTokenTTL time.Duration
}

// Deps bundles the auth service dependencies. A struct keeps the constructor
// within the ≤4-parameter limit while staying explicit about wiring.
//
// Tx + Orgs power tenant signup: Register creates the user, their personal
// organization, and the owner membership in ONE transaction so a user is never
// persisted without a home organization. Orgs is also used on login to resolve
// the active organization stamped into the access token.
type Deps struct {
	Users           repository.UserRepository
	Orgs            repository.OrganizationRepository
	Resets          repository.PasswordResetRepository
	Verifications   repository.EmailVerificationRepository
	RefreshTokens   repository.RefreshTokenRepository
	Tx              repository.TxManager
	Email           service.EmailService
	JWTMaker        *jwt.Maker
	FrontendURL     string
	RefreshTokenTTL time.Duration
}

func NewService(deps Deps) *Service {
	return &Service{
		users:           deps.Users,
		orgs:            deps.Orgs,
		resets:          deps.Resets,
		verifications:   deps.Verifications,
		refreshTokens:   deps.RefreshTokens,
		tx:              deps.Tx,
		email:           deps.Email,
		jwtMaker:        deps.JWTMaker,
		frontendURL:     deps.FrontendURL,
		refreshTokenTTL: deps.RefreshTokenTTL,
	}
}

func (s *Service) Register(ctx context.Context, email, name, password string) (*user.User, string, string, error) {
	_, err := s.users.FindByEmail(ctx, email)
	if err == nil {
		return nil, "", "", domain.ErrAlreadyExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, "", "", err
	}

	hashed, err := hash.Password(password)
	if err != nil {
		return nil, "", "", err
	}

	u, err := user.New(email, name, hashed)
	if err != nil {
		return nil, "", "", err
	}

	// Create the user, their personal organization, and the owner membership in ONE
	// transaction. activeOrgID is set inside the callback once the org has an ID, so
	// the access token can be stamped with it below.
	var activeOrgID string
	err = s.tx.WithSignupTx(ctx, func(users repository.UserRepository, orgs repository.OrganizationRepository, members repository.OrganizationMemberRepository) error {
		if err := users.Create(ctx, u); err != nil {
			return err
		}
		o, err := org.New(name, "", u.ID)
		if err != nil {
			return err
		}
		// Personal-org slugs must be unique; suffix with a short slice of the user
		// id so two users named the same do not collide.
		o.Slug = uniqueSlug(o.Slug, u.ID)
		if err := orgs.Create(ctx, o); err != nil {
			return fmt.Errorf("creating personal organization: %w", err)
		}
		owner, err := org.NewMember(o.ID, u.ID, org.RoleOwner)
		if err != nil {
			return err
		}
		if err := members.Add(ctx, owner); err != nil {
			return fmt.Errorf("adding owner membership: %w", err)
		}
		activeOrgID = o.ID
		return nil
	})
	if err != nil {
		return nil, "", "", err
	}

	// Send verification email asynchronously (don't block registration)
	if s.email != nil && s.verifications != nil {
		_ = s.sendVerificationEmail(ctx, u)
	}

	access, refresh, err := s.issueTokens(ctx, u, activeOrgID)
	if err != nil {
		return nil, "", "", err
	}

	return u, access, refresh, nil
}

// uniqueSlug appends a short, stable suffix derived from the user id to a base
// slug so per-user personal organizations never collide on the unique slug index.
func uniqueSlug(base, userID string) string {
	suffix := userID
	if len(suffix) > 8 {
		suffix = suffix[:8]
	}
	if base == "" {
		base = "org"
	}
	return base + "-" + suffix
}

func (s *Service) Login(ctx context.Context, email, password string) (*user.User, string, string, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", "", domain.ErrUnauthorized
	}

	if !hash.Check(password, u.PasswordHash) {
		return nil, "", "", domain.ErrUnauthorized
	}

	access, refresh, err := s.issueTokens(ctx, u, s.activeOrgID(ctx, u.ID))
	if err != nil {
		return nil, "", "", err
	}

	return u, access, refresh, nil
}

func (s *Service) OAuthCallback(ctx context.Context, oauthUser *service.OAuthUser) (*user.User, string, string, error) {
	u, err := s.users.FindByProvider(ctx, oauthUser.Provider, oauthUser.ProviderID)
	if errors.Is(err, domain.ErrNotFound) {
		u = user.NewOAuth(oauthUser.Email, oauthUser.Name, oauthUser.Provider, oauthUser.ProviderID)
		u.AvatarURL = oauthUser.AvatarURL
		if err := s.users.Create(ctx, u); err != nil {
			return nil, "", "", err
		}
	} else if err != nil {
		return nil, "", "", err
	}

	// An OAuth user signing in for the first time has no organization yet — ensure
	// one exists so they land in a tenant just like an email signup.
	orgID, err := s.ensureOrg(ctx, u)
	if err != nil {
		return nil, "", "", err
	}

	access, refresh, err := s.issueTokens(ctx, u, orgID)
	if err != nil {
		return nil, "", "", err
	}

	return u, access, refresh, nil
}

// issueTokens mints a short-lived access token (carrying the active org) and a
// long-lived opaque refresh token, persisting only the SHA-256 hash of the refresh
// token. The raw refresh token is returned to the caller and never stored.
func (s *Service) issueTokens(ctx context.Context, u *user.User, orgID string) (access, refresh string, err error) {
	access, err = s.jwtMaker.GenerateWithOrg(u.ID, string(u.Role), orgID)
	if err != nil {
		return "", "", err
	}

	refresh, err = s.storeRefreshToken(ctx, u.ID)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

// activeOrgID returns the user's default/personal organization id, or "" if it
// cannot be resolved. Resolving the org is best-effort here: a token without an
// org claim still works because the middleware falls back to a default-org lookup.
func (s *Service) activeOrgID(ctx context.Context, userID string) string {
	if s.orgs == nil {
		return ""
	}
	o, err := s.orgs.FindDefaultForUser(ctx, userID)
	if err != nil {
		return ""
	}
	return o.ID
}

// ensureOrg returns the user's default organization id, creating a personal
// organization (and owner membership) on the fly if none exists yet — used by the
// OAuth path where there is no explicit registration step.
func (s *Service) ensureOrg(ctx context.Context, u *user.User) (string, error) {
	if s.orgs == nil {
		return "", nil
	}
	if o, err := s.orgs.FindDefaultForUser(ctx, u.ID); err == nil {
		return o.ID, nil
	} else if !errors.Is(err, domain.ErrNotFound) {
		return "", err
	}

	var orgID string
	err := s.tx.WithSignupTx(ctx, func(_ repository.UserRepository, orgs repository.OrganizationRepository, members repository.OrganizationMemberRepository) error {
		o, err := org.New(u.Name, "", u.ID)
		if err != nil {
			return err
		}
		o.Slug = uniqueSlug(o.Slug, u.ID)
		if err := orgs.Create(ctx, o); err != nil {
			return fmt.Errorf("creating personal organization: %w", err)
		}
		owner, err := org.NewMember(o.ID, u.ID, org.RoleOwner)
		if err != nil {
			return err
		}
		if err := members.Add(ctx, owner); err != nil {
			return fmt.Errorf("adding owner membership: %w", err)
		}
		orgID = o.ID
		return nil
	})
	return orgID, err
}

// storeRefreshToken generates an opaque refresh token, persists its hash, and
// returns the raw token to hand back to the client.
func (s *Service) storeRefreshToken(ctx context.Context, userID string) (string, error) {
	raw, err := jwt.GenerateRefreshToken()
	if err != nil {
		return "", fmt.Errorf("generating refresh token: %w", err)
	}

	rt := &repository.RefreshToken{
		UserID:    userID,
		TokenHash: jwt.HashRefreshToken(raw),
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
	}
	if err := s.refreshTokens.Create(ctx, rt); err != nil {
		return "", fmt.Errorf("storing refresh token: %w", err)
	}

	return raw, nil
}

// Refresh validates an opaque refresh token, rotates it (revokes the presented
// token and issues a fresh one), and returns a new access token. Rotation makes
// refresh-token reuse detectable: a revoked/expired token is rejected outright.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (newAccess, newRefresh string, u *user.User, err error) {
	hashed := jwt.HashRefreshToken(refreshToken)

	rt, err := s.refreshTokens.FindByHash(ctx, hashed)
	if err != nil {
		// Missing token (or any lookup failure) is unauthorized — do not leak detail.
		return "", "", nil, domain.ErrUnauthorized
	}

	if !rt.IsValid(time.Now()) {
		// Expired or already-revoked: reject. If it was revoked, this is a
		// potential reuse signal a caller could escalate in the future.
		return "", "", nil, domain.ErrUnauthorized
	}

	u, err = s.users.FindByID(ctx, rt.UserID)
	if err != nil {
		return "", "", nil, domain.ErrUnauthorized
	}

	// Rotate: revoke the presented token, then mint a replacement pair.
	if err := s.refreshTokens.Revoke(ctx, hashed); err != nil {
		return "", "", nil, fmt.Errorf("revoking rotated refresh token: %w", err)
	}

	newAccess, newRefresh, err = s.issueTokens(ctx, u, s.activeOrgID(ctx, u.ID))
	if err != nil {
		return "", "", nil, err
	}

	return newAccess, newRefresh, u, nil
}

// Logout revokes a single refresh token so it can no longer be used to obtain
// new access tokens. The access token remains valid until it expires (short TTL
// is the accepted revocation window).
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	hashed := jwt.HashRefreshToken(refreshToken)
	if err := s.refreshTokens.Revoke(ctx, hashed); err != nil {
		// An unknown/already-revoked token is not an error for logout — the
		// desired end state (no usable session) is already met.
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("revoking refresh token on logout: %w", err)
	}
	return nil
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

	// Invalidate every existing session: a password reset must log out all
	// refresh tokens so a compromised/old session cannot survive the reset.
	if s.refreshTokens != nil {
		if err := s.refreshTokens.RevokeAllForUser(ctx, u.ID); err != nil {
			return fmt.Errorf("revoking refresh tokens after password reset: %w", err)
		}
	}

	return nil
}
