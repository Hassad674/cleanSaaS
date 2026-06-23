package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainorg "github.com/hassad/boilerplateSaaS/backend/internal/domain/org"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
	"github.com/hassad/boilerplateSaaS/backend/pkg/hash"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
)

// Mock implementations

type mockUserRepo struct {
	findByEmailFn func(ctx context.Context, email string) (*user.User, error)
	findByIDFn    func(ctx context.Context, id string) (*user.User, error)
	createFn      func(ctx context.Context, u *user.User) error
	updateFn      func(ctx context.Context, u *user.User) error
}

func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error {
	if m.createFn != nil {
		return m.createFn(ctx, u)
	}
	return nil
}
func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*user.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email)
	}
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) FindByProvider(_ context.Context, _, _ string) (*user.User, error) {
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) FindByStripeID(_ context.Context, _ string) (*user.User, error) {
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) Update(ctx context.Context, u *user.User) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, u)
	}
	return nil
}
func (m *mockUserRepo) Delete(_ context.Context, _ string) error { return nil }
func (m *mockUserRepo) List(_ context.Context, _, _ int) ([]*user.User, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) Search(_ context.Context, _ string, _, _ int) ([]*user.User, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) Count(_ context.Context) (int, error) { return 0, nil }

type mockResetRepo struct {
	createFn      func(ctx context.Context, pr *repository.PasswordReset) error
	findByTokenFn func(ctx context.Context, token string) (*repository.PasswordReset, error)
	markUsedFn    func(ctx context.Context, id string) error
}

func (m *mockResetRepo) Create(ctx context.Context, pr *repository.PasswordReset) error {
	if m.createFn != nil {
		return m.createFn(ctx, pr)
	}
	return nil
}
func (m *mockResetRepo) FindByToken(ctx context.Context, token string) (*repository.PasswordReset, error) {
	if m.findByTokenFn != nil {
		return m.findByTokenFn(ctx, token)
	}
	return nil, domain.ErrNotFound
}
func (m *mockResetRepo) MarkUsed(ctx context.Context, id string) error {
	if m.markUsedFn != nil {
		return m.markUsedFn(ctx, id)
	}
	return nil
}
func (m *mockResetRepo) DeleteExpired(_ context.Context) error { return nil }

type mockVerificationRepo struct {
	createFn         func(ctx context.Context, ev *repository.EmailVerification) error
	findByTokenFn    func(ctx context.Context, token string) (*repository.EmailVerification, error)
	deleteByUserIDFn func(ctx context.Context, userID string) error
}

func (m *mockVerificationRepo) Create(ctx context.Context, ev *repository.EmailVerification) error {
	if m.createFn != nil {
		return m.createFn(ctx, ev)
	}
	return nil
}
func (m *mockVerificationRepo) FindByToken(ctx context.Context, token string) (*repository.EmailVerification, error) {
	if m.findByTokenFn != nil {
		return m.findByTokenFn(ctx, token)
	}
	return nil, domain.ErrNotFound
}
func (m *mockVerificationRepo) DeleteByUserID(ctx context.Context, userID string) error {
	if m.deleteByUserIDFn != nil {
		return m.deleteByUserIDFn(ctx, userID)
	}
	return nil
}
func (m *mockVerificationRepo) DeleteExpired(_ context.Context) error { return nil }

type mockRefreshRepo struct {
	createFn           func(ctx context.Context, token *repository.RefreshToken) error
	findByHashFn       func(ctx context.Context, hash string) (*repository.RefreshToken, error)
	revokeFn           func(ctx context.Context, hash string) error
	revokeAllForUserFn func(ctx context.Context, userID string) error
}

func (m *mockRefreshRepo) Create(ctx context.Context, token *repository.RefreshToken) error {
	if m.createFn != nil {
		return m.createFn(ctx, token)
	}
	return nil
}
func (m *mockRefreshRepo) FindByHash(ctx context.Context, hash string) (*repository.RefreshToken, error) {
	if m.findByHashFn != nil {
		return m.findByHashFn(ctx, hash)
	}
	return nil, domain.ErrNotFound
}
func (m *mockRefreshRepo) Revoke(ctx context.Context, hash string) error {
	if m.revokeFn != nil {
		return m.revokeFn(ctx, hash)
	}
	return nil
}
func (m *mockRefreshRepo) RevokeAllForUser(ctx context.Context, userID string) error {
	if m.revokeAllForUserFn != nil {
		return m.revokeAllForUserFn(ctx, userID)
	}
	return nil
}
func (m *mockRefreshRepo) DeleteExpired(_ context.Context) error { return nil }

type mockEmailSvc struct {
	sendFn         func(ctx context.Context, email service.Email) error
	sendTemplateFn func(ctx context.Context, to, template string, data map[string]string) error
}

func (m *mockEmailSvc) Send(ctx context.Context, email service.Email) error {
	if m.sendFn != nil {
		return m.sendFn(ctx, email)
	}
	return nil
}
func (m *mockEmailSvc) SendTemplate(ctx context.Context, to, template string, data map[string]string) error {
	if m.sendTemplateFn != nil {
		return m.sendTemplateFn(ctx, to, template, data)
	}
	return nil
}

// testDeps bundles the optional mocks a test wants to customize. Nil fields are
// filled with zero-value mocks so each test only sets what it cares about.
type testDeps struct {
	userRepo    *mockUserRepo
	resetRepo   *mockResetRepo
	verifyRepo  *mockVerificationRepo
	refreshRepo *mockRefreshRepo
	emailSvc    *mockEmailSvc
}

// mockOrgRepo is a minimal OrganizationRepository for auth tests. FindDefaultForUser
// reports "not found" so the OAuth path creates a personal org; Create assigns an id.
type mockOrgRepo struct{}

func (m *mockOrgRepo) Create(_ context.Context, o *domainorg.Organization) error {
	o.ID = "org-1"
	return nil
}
func (m *mockOrgRepo) FindByID(_ context.Context, _ string) (*domainorg.Organization, error) {
	return nil, domain.ErrNotFound
}
func (m *mockOrgRepo) FindBySlug(_ context.Context, _ string) (*domainorg.Organization, error) {
	return nil, domain.ErrNotFound
}
func (m *mockOrgRepo) FindDefaultForUser(_ context.Context, _ string) (*domainorg.Organization, error) {
	return nil, domain.ErrNotFound
}

type mockOrgMemberRepo struct{}

func (m *mockOrgMemberRepo) Add(_ context.Context, _ *domainorg.Member) error { return nil }
func (m *mockOrgMemberRepo) FindByOrgAndUser(_ context.Context, _, _ string) (*domainorg.Member, error) {
	return nil, domain.ErrNotFound
}
func (m *mockOrgMemberRepo) IsMember(_ context.Context, _, _ string) (bool, error) { return true, nil }

// mockTxManager runs the signup callback directly with the test's user repo and
// stub org repos — no real transaction is needed in a unit test.
type mockTxManager struct{ users repository.UserRepository }

func (m *mockTxManager) WithTeamTx(ctx context.Context, fn func(teams repository.TeamRepository, members repository.TeamMemberRepository) error) error {
	return fn(nil, nil)
}
func (m *mockTxManager) WithSignupTx(ctx context.Context, fn func(users repository.UserRepository, orgs repository.OrganizationRepository, members repository.OrganizationMemberRepository) error) error {
	return fn(m.users, &mockOrgRepo{}, &mockOrgMemberRepo{})
}

func newTestService(d testDeps) *Service {
	if d.userRepo == nil {
		d.userRepo = &mockUserRepo{}
	}
	if d.resetRepo == nil {
		d.resetRepo = &mockResetRepo{}
	}
	if d.verifyRepo == nil {
		d.verifyRepo = &mockVerificationRepo{}
	}
	if d.refreshRepo == nil {
		d.refreshRepo = &mockRefreshRepo{}
	}
	deps := Deps{
		Users:           d.userRepo,
		Orgs:            &mockOrgRepo{},
		Resets:          d.resetRepo,
		Verifications:   d.verifyRepo,
		RefreshTokens:   d.refreshRepo,
		Tx:              &mockTxManager{users: d.userRepo},
		JWTMaker:        jwt.NewMaker("secret"),
		FrontendURL:     "http://localhost:3006",
		RefreshTokenTTL: 720 * time.Hour,
	}
	// Assign Email only when set, so a nil mock yields a true-nil interface
	// (a typed nil *mockEmailSvc would make s.email != nil checks pass).
	if d.emailSvc != nil {
		deps.Email = d.emailSvc
	}
	return NewService(deps)
}

// === ForgotPassword tests ===

func TestAuthService_ForgotPassword_UserExists(t *testing.T) {
	var sentEmail bool
	userRepo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", Email: "test@test.com", Name: "Test"}, nil
		},
	}
	resetRepo := &mockResetRepo{
		createFn: func(_ context.Context, _ *repository.PasswordReset) error { return nil },
	}
	emailSvc := &mockEmailSvc{
		sendTemplateFn: func(_ context.Context, _, _ string, _ map[string]string) error {
			sentEmail = true
			return nil
		},
	}

	svc := newTestService(testDeps{userRepo: userRepo, resetRepo: resetRepo, emailSvc: emailSvc})
	err := svc.ForgotPassword(context.Background(), "test@test.com", "")
	assert.NoError(t, err)
	assert.True(t, sentEmail, "email should have been sent")
}

func TestAuthService_ForgotPassword_UserNotFound(t *testing.T) {
	userRepo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*user.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := newTestService(testDeps{userRepo: userRepo, resetRepo: &mockResetRepo{}})
	err := svc.ForgotPassword(context.Background(), "noone@test.com", "")
	assert.NoError(t, err, "should not return error for non-existent user")
}

// === ResetPassword tests ===

func TestAuthService_ResetPassword_ValidToken(t *testing.T) {
	var updated bool
	hashed, _ := hash.Password("oldpass")

	userRepo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", PasswordHash: hashed}, nil
		},
		updateFn: func(_ context.Context, _ *user.User) error {
			updated = true
			return nil
		},
	}
	resetRepo := &mockResetRepo{
		findByTokenFn: func(_ context.Context, _ string) (*repository.PasswordReset, error) {
			return &repository.PasswordReset{
				ID:        "r1",
				UserID:    "u1",
				Token:     "valid-token",
				ExpiresAt: time.Now().Add(time.Hour),
				Used:      false,
			}, nil
		},
		markUsedFn: func(_ context.Context, _ string) error { return nil },
	}

	svc := newTestService(testDeps{userRepo: userRepo, resetRepo: resetRepo})
	err := svc.ResetPassword(context.Background(), "valid-token", "newpassword")
	assert.NoError(t, err)
	assert.True(t, updated, "user password should have been updated")
}

func TestAuthService_ResetPassword_ExpiredToken(t *testing.T) {
	resetRepo := &mockResetRepo{
		findByTokenFn: func(_ context.Context, _ string) (*repository.PasswordReset, error) {
			return &repository.PasswordReset{
				ID:        "r1",
				UserID:    "u1",
				Token:     "expired-token",
				ExpiresAt: time.Now().Add(-time.Hour),
				Used:      false,
			}, nil
		},
	}

	svc := newTestService(testDeps{resetRepo: resetRepo})
	err := svc.ResetPassword(context.Background(), "expired-token", "newpassword")
	assert.ErrorIs(t, err, domain.ErrExpiredToken)
}

func TestAuthService_ResetPassword_UsedToken(t *testing.T) {
	resetRepo := &mockResetRepo{
		findByTokenFn: func(_ context.Context, _ string) (*repository.PasswordReset, error) {
			return &repository.PasswordReset{
				ID:        "r1",
				UserID:    "u1",
				Token:     "used-token",
				ExpiresAt: time.Now().Add(time.Hour),
				Used:      true,
			}, nil
		},
	}

	svc := newTestService(testDeps{resetRepo: resetRepo})
	err := svc.ResetPassword(context.Background(), "used-token", "newpassword")
	assert.ErrorIs(t, err, domain.ErrInvalidToken)
}

// === VerifyEmail tests ===

func TestAuthService_VerifyEmail_ValidToken(t *testing.T) {
	var updatedUser *user.User
	var deletedUserID string

	userRepo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", Email: "test@test.com", EmailVerified: false}, nil
		},
		updateFn: func(_ context.Context, u *user.User) error {
			updatedUser = u
			return nil
		},
	}
	verifyRepo := &mockVerificationRepo{
		findByTokenFn: func(_ context.Context, _ string) (*repository.EmailVerification, error) {
			return &repository.EmailVerification{
				ID:        "v1",
				UserID:    "u1",
				Token:     "valid-token",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}, nil
		},
		deleteByUserIDFn: func(_ context.Context, userID string) error {
			deletedUserID = userID
			return nil
		},
	}

	svc := newTestService(testDeps{userRepo: userRepo, verifyRepo: verifyRepo})
	err := svc.VerifyEmail(context.Background(), "valid-token")
	assert.NoError(t, err)
	assert.True(t, updatedUser.EmailVerified, "user should be marked as verified")
	assert.Equal(t, "u1", deletedUserID, "verification tokens should be cleaned up")
}

func TestAuthService_VerifyEmail_ExpiredToken(t *testing.T) {
	verifyRepo := &mockVerificationRepo{
		findByTokenFn: func(_ context.Context, _ string) (*repository.EmailVerification, error) {
			return &repository.EmailVerification{
				ID:        "v1",
				UserID:    "u1",
				Token:     "expired-token",
				ExpiresAt: time.Now().Add(-time.Hour),
			}, nil
		},
	}

	svc := newTestService(testDeps{verifyRepo: verifyRepo})
	err := svc.VerifyEmail(context.Background(), "expired-token")
	assert.ErrorIs(t, err, domain.ErrExpiredToken)
}

func TestAuthService_VerifyEmail_AlreadyVerified(t *testing.T) {
	userRepo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", Email: "test@test.com", EmailVerified: true}, nil
		},
	}
	verifyRepo := &mockVerificationRepo{
		findByTokenFn: func(_ context.Context, _ string) (*repository.EmailVerification, error) {
			return &repository.EmailVerification{
				ID:        "v1",
				UserID:    "u1",
				Token:     "valid-token",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}, nil
		},
	}

	svc := newTestService(testDeps{userRepo: userRepo, verifyRepo: verifyRepo})
	err := svc.VerifyEmail(context.Background(), "valid-token")
	assert.NoError(t, err, "should succeed for already verified user")
}

func TestAuthService_VerifyEmail_InvalidToken(t *testing.T) {
	verifyRepo := &mockVerificationRepo{
		findByTokenFn: func(_ context.Context, _ string) (*repository.EmailVerification, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := newTestService(testDeps{verifyRepo: verifyRepo})
	err := svc.VerifyEmail(context.Background(), "bogus-token")
	assert.ErrorIs(t, err, domain.ErrInvalidToken)
}

// === ResendVerification tests ===

func TestAuthService_ResendVerification_Success(t *testing.T) {
	var sentTemplate string
	userRepo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", Email: "test@test.com", Name: "Test", EmailVerified: false}, nil
		},
	}
	verifyRepo := &mockVerificationRepo{}
	emailSvc := &mockEmailSvc{
		sendTemplateFn: func(_ context.Context, _, template string, _ map[string]string) error {
			sentTemplate = template
			return nil
		},
	}

	svc := newTestService(testDeps{userRepo: userRepo, verifyRepo: verifyRepo, emailSvc: emailSvc})
	err := svc.ResendVerification(context.Background(), "u1")
	assert.NoError(t, err)
	assert.Equal(t, "verification", sentTemplate, "should send verification template")
}

func TestAuthService_ResendVerification_AlreadyVerified(t *testing.T) {
	userRepo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", Email: "test@test.com", EmailVerified: true}, nil
		},
	}

	svc := newTestService(testDeps{userRepo: userRepo})
	err := svc.ResendVerification(context.Background(), "u1")
	assert.NoError(t, err, "should be a no-op for already verified users")
}

// === Login / Register issue refresh tokens ===

func TestAuthService_Login_ReturnsAccessAndRefreshTokens(t *testing.T) {
	hashed, _ := hash.Password("correct-password")
	var storedRefresh *repository.RefreshToken

	userRepo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", Email: "test@test.com", Role: user.RoleMember, PasswordHash: hashed}, nil
		},
	}
	refreshRepo := &mockRefreshRepo{
		createFn: func(_ context.Context, rt *repository.RefreshToken) error {
			storedRefresh = rt
			return nil
		},
	}

	svc := newTestService(testDeps{userRepo: userRepo, refreshRepo: refreshRepo})
	u, access, refresh, err := svc.Login(context.Background(), "test@test.com", "correct-password")
	assert.NoError(t, err)
	assert.Equal(t, "u1", u.ID)
	assert.NotEmpty(t, access, "access token should be issued")
	assert.NotEmpty(t, refresh, "refresh token should be issued")
	assert.NotEqual(t, access, refresh, "access and refresh tokens must differ")

	// Only the hash is persisted — never the raw refresh token.
	assert.NotNil(t, storedRefresh)
	assert.Equal(t, jwt.HashRefreshToken(refresh), storedRefresh.TokenHash)
	assert.NotEqual(t, refresh, storedRefresh.TokenHash, "raw refresh token must not be stored")
	assert.Equal(t, "u1", storedRefresh.UserID)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	hashed, _ := hash.Password("correct-password")
	userRepo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", PasswordHash: hashed, Role: user.RoleMember}, nil
		},
	}

	svc := newTestService(testDeps{userRepo: userRepo})
	_, _, _, err := svc.Login(context.Background(), "test@test.com", "wrong-password")
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestAuthService_Register_ReturnsRefreshToken(t *testing.T) {
	var createdRefresh bool
	userRepo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*user.User, error) {
			return nil, domain.ErrNotFound
		},
		createFn: func(_ context.Context, u *user.User) error {
			u.ID = "u1"
			return nil
		},
	}
	refreshRepo := &mockRefreshRepo{
		createFn: func(_ context.Context, _ *repository.RefreshToken) error {
			createdRefresh = true
			return nil
		},
	}

	svc := newTestService(testDeps{userRepo: userRepo, refreshRepo: refreshRepo})
	_, access, refresh, err := svc.Register(context.Background(), "new@test.com", "New User", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)
	assert.True(t, createdRefresh, "register should persist a refresh token")
}

// === Refresh (rotation) tests ===

func TestAuthService_Refresh_RotatesToken(t *testing.T) {
	const rawToken = "the-presented-refresh-token"
	var revokedHash string
	var newStoredHash string

	userRepo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", Role: user.RoleMember}, nil
		},
	}
	refreshRepo := &mockRefreshRepo{
		findByHashFn: func(_ context.Context, hash string) (*repository.RefreshToken, error) {
			return &repository.RefreshToken{
				ID:        "rt1",
				UserID:    "u1",
				TokenHash: hash,
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}, nil
		},
		revokeFn: func(_ context.Context, hash string) error {
			revokedHash = hash
			return nil
		},
		createFn: func(_ context.Context, rt *repository.RefreshToken) error {
			newStoredHash = rt.TokenHash
			return nil
		},
	}

	svc := newTestService(testDeps{userRepo: userRepo, refreshRepo: refreshRepo})
	newAccess, newRefresh, u, err := svc.Refresh(context.Background(), rawToken)
	assert.NoError(t, err)
	assert.Equal(t, "u1", u.ID)
	assert.NotEmpty(t, newAccess)
	assert.NotEmpty(t, newRefresh)

	// The presented token is revoked (rotation) ...
	assert.Equal(t, jwt.HashRefreshToken(rawToken), revokedHash, "presented token must be revoked")
	// ... and a brand-new, different refresh token is stored.
	assert.NotEqual(t, rawToken, newRefresh)
	assert.Equal(t, jwt.HashRefreshToken(newRefresh), newStoredHash)
	assert.NotEqual(t, revokedHash, newStoredHash, "rotation must produce a new token")
}

func TestAuthService_Refresh_RejectsMissingToken(t *testing.T) {
	refreshRepo := &mockRefreshRepo{
		findByHashFn: func(_ context.Context, _ string) (*repository.RefreshToken, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := newTestService(testDeps{refreshRepo: refreshRepo})
	_, _, _, err := svc.Refresh(context.Background(), "unknown-token")
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestAuthService_Refresh_RejectsRevokedToken(t *testing.T) {
	revokedAt := time.Now().Add(-time.Minute)
	refreshRepo := &mockRefreshRepo{
		findByHashFn: func(_ context.Context, hash string) (*repository.RefreshToken, error) {
			return &repository.RefreshToken{
				ID:        "rt1",
				UserID:    "u1",
				TokenHash: hash,
				ExpiresAt: time.Now().Add(24 * time.Hour),
				RevokedAt: &revokedAt,
			}, nil
		},
	}

	svc := newTestService(testDeps{refreshRepo: refreshRepo})
	_, _, _, err := svc.Refresh(context.Background(), "revoked-token")
	assert.ErrorIs(t, err, domain.ErrUnauthorized, "revoked token must be rejected")
}

func TestAuthService_Refresh_RejectsExpiredToken(t *testing.T) {
	refreshRepo := &mockRefreshRepo{
		findByHashFn: func(_ context.Context, hash string) (*repository.RefreshToken, error) {
			return &repository.RefreshToken{
				ID:        "rt1",
				UserID:    "u1",
				TokenHash: hash,
				ExpiresAt: time.Now().Add(-time.Hour),
			}, nil
		},
	}

	svc := newTestService(testDeps{refreshRepo: refreshRepo})
	_, _, _, err := svc.Refresh(context.Background(), "expired-token")
	assert.ErrorIs(t, err, domain.ErrUnauthorized, "expired token must be rejected")
}

// === Logout tests ===

func TestAuthService_Logout_RevokesToken(t *testing.T) {
	const rawToken = "session-refresh-token"
	var revokedHash string

	refreshRepo := &mockRefreshRepo{
		revokeFn: func(_ context.Context, hash string) error {
			revokedHash = hash
			return nil
		},
	}

	svc := newTestService(testDeps{refreshRepo: refreshRepo})
	err := svc.Logout(context.Background(), rawToken)
	assert.NoError(t, err)
	assert.Equal(t, jwt.HashRefreshToken(rawToken), revokedHash, "logout must revoke the presented token")
}

func TestAuthService_Logout_UnknownTokenIsNoOp(t *testing.T) {
	refreshRepo := &mockRefreshRepo{
		revokeFn: func(_ context.Context, _ string) error {
			return domain.ErrNotFound
		},
	}

	svc := newTestService(testDeps{refreshRepo: refreshRepo})
	err := svc.Logout(context.Background(), "already-gone")
	assert.NoError(t, err, "logging out an unknown token should be idempotent")
}

// === Password reset revokes all sessions ===

func TestAuthService_ResetPassword_RevokesAllRefreshTokens(t *testing.T) {
	hashed, _ := hash.Password("oldpass")
	var revokedUserID string

	userRepo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", PasswordHash: hashed}, nil
		},
	}
	resetRepo := &mockResetRepo{
		findByTokenFn: func(_ context.Context, _ string) (*repository.PasswordReset, error) {
			return &repository.PasswordReset{
				ID:        "r1",
				UserID:    "u1",
				Token:     "valid-token",
				ExpiresAt: time.Now().Add(time.Hour),
			}, nil
		},
	}
	refreshRepo := &mockRefreshRepo{
		revokeAllForUserFn: func(_ context.Context, userID string) error {
			revokedUserID = userID
			return nil
		},
	}

	svc := newTestService(testDeps{userRepo: userRepo, resetRepo: resetRepo, refreshRepo: refreshRepo})
	err := svc.ResetPassword(context.Background(), "valid-token", "newpassword")
	assert.NoError(t, err)
	assert.Equal(t, "u1", revokedUserID, "password reset must revoke all of the user's refresh tokens")
}
