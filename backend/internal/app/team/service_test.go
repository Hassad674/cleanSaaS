package team

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainteam "github.com/hassad/boilerplateSaaS/backend/internal/domain/team"
)

// --- Mocks ---

type mockTeamRepo struct {
	createFn     func(ctx context.Context, t *domainteam.Team) error
	findByIDFn   func(ctx context.Context, id string) (*domainteam.Team, error)
	findBySlugFn func(ctx context.Context, slug string) (*domainteam.Team, error)
	updateFn     func(ctx context.Context, t *domainteam.Team) error
	deleteFn     func(ctx context.Context, id string) error
	listByUserFn func(ctx context.Context, userID string) ([]*domainteam.Team, error)
}

func (m *mockTeamRepo) Create(ctx context.Context, t *domainteam.Team) error {
	if m.createFn != nil {
		return m.createFn(ctx, t)
	}
	t.ID = "team-1"
	return nil
}

func (m *mockTeamRepo) FindByID(ctx context.Context, id string) (*domainteam.Team, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *mockTeamRepo) FindBySlug(ctx context.Context, slug string) (*domainteam.Team, error) {
	if m.findBySlugFn != nil {
		return m.findBySlugFn(ctx, slug)
	}
	return nil, domain.ErrNotFound
}

func (m *mockTeamRepo) Update(ctx context.Context, t *domainteam.Team) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, t)
	}
	return nil
}

func (m *mockTeamRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockTeamRepo) ListByUserID(ctx context.Context, userID string) ([]*domainteam.Team, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID)
	}
	return nil, nil
}

type mockMemberRepo struct {
	addFn               func(ctx context.Context, member *domainteam.TeamMember) error
	findByIDFn          func(ctx context.Context, id string) (*domainteam.TeamMember, error)
	findByTeamAndUserFn func(ctx context.Context, teamID, userID string) (*domainteam.TeamMember, error)
	findByInviteTokenFn func(ctx context.Context, token string) (*domainteam.TeamMember, error)
	updateFn            func(ctx context.Context, member *domainteam.TeamMember) error
	removeFn            func(ctx context.Context, teamID, userID string) error
	listByTeamIDFn      func(ctx context.Context, teamID string, offset, limit int) ([]*domainteam.TeamMember, int, error)
	countByTeamIDFn     func(ctx context.Context, teamID string) (int, error)
}

func (m *mockMemberRepo) Add(ctx context.Context, member *domainteam.TeamMember) error {
	if m.addFn != nil {
		return m.addFn(ctx, member)
	}
	member.ID = "member-1"
	return nil
}

func (m *mockMemberRepo) FindByID(ctx context.Context, id string) (*domainteam.TeamMember, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *mockMemberRepo) FindByTeamAndUser(ctx context.Context, teamID, userID string) (*domainteam.TeamMember, error) {
	if m.findByTeamAndUserFn != nil {
		return m.findByTeamAndUserFn(ctx, teamID, userID)
	}
	return nil, domain.ErrNotFound
}

func (m *mockMemberRepo) FindByInviteToken(ctx context.Context, token string) (*domainteam.TeamMember, error) {
	if m.findByInviteTokenFn != nil {
		return m.findByInviteTokenFn(ctx, token)
	}
	return nil, domain.ErrNotFound
}

func (m *mockMemberRepo) Update(ctx context.Context, member *domainteam.TeamMember) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, member)
	}
	return nil
}

func (m *mockMemberRepo) Remove(ctx context.Context, teamID, userID string) error {
	if m.removeFn != nil {
		return m.removeFn(ctx, teamID, userID)
	}
	return nil
}

func (m *mockMemberRepo) ListByTeamID(ctx context.Context, teamID string, offset, limit int) ([]*domainteam.TeamMember, int, error) {
	if m.listByTeamIDFn != nil {
		return m.listByTeamIDFn(ctx, teamID, offset, limit)
	}
	return nil, 0, nil
}

func (m *mockMemberRepo) CountByTeamID(ctx context.Context, teamID string) (int, error) {
	if m.countByTeamIDFn != nil {
		return m.countByTeamIDFn(ctx, teamID)
	}
	return 0, nil
}

// --- CreateTeam ---

func TestService_CreateTeam_Success(t *testing.T) {
	teamRepo := &mockTeamRepo{}
	memberRepo := &mockMemberRepo{}
	svc := NewService(teamRepo, memberRepo)

	tm, err := svc.CreateTeam(context.Background(), "user-1", "My Team")
	assert.NoError(t, err)
	assert.Equal(t, "team-1", tm.ID)
	assert.Equal(t, "My Team", tm.Name)
	assert.Equal(t, "my-team", tm.Slug)
	assert.Equal(t, "user-1", tm.OwnerID)
}

func TestService_CreateTeam_ValidationError(t *testing.T) {
	teamRepo := &mockTeamRepo{}
	memberRepo := &mockMemberRepo{}
	svc := NewService(teamRepo, memberRepo)

	_, err := svc.CreateTeam(context.Background(), "user-1", "A") // too short
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestService_CreateTeam_EmptyName(t *testing.T) {
	teamRepo := &mockTeamRepo{}
	memberRepo := &mockMemberRepo{}
	svc := NewService(teamRepo, memberRepo)

	_, err := svc.CreateTeam(context.Background(), "user-1", "")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestService_CreateTeam_AddsOwnerAsMember(t *testing.T) {
	var addedMember *domainteam.TeamMember
	teamRepo := &mockTeamRepo{}
	memberRepo := &mockMemberRepo{
		addFn: func(_ context.Context, member *domainteam.TeamMember) error {
			addedMember = member
			member.ID = "member-1"
			return nil
		},
	}
	svc := NewService(teamRepo, memberRepo)

	_, err := svc.CreateTeam(context.Background(), "user-1", "My Team")
	assert.NoError(t, err)
	assert.NotNil(t, addedMember)
	assert.Equal(t, "user-1", addedMember.UserID)
	assert.Equal(t, domainteam.RoleOwner, addedMember.Role)
	assert.Equal(t, domainteam.InviteAccepted, addedMember.InviteStatus)
	assert.NotNil(t, addedMember.JoinedAt)
}

// --- GetTeam ---

// memberRepoWithMember returns a member repo whose membership lookup always succeeds.
func memberRepoWithMember(role domainteam.Role) *mockMemberRepo {
	return &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: role}, nil
		},
	}
}

func TestService_GetTeam_Success(t *testing.T) {
	teamRepo := &mockTeamRepo{
		findByIDFn: func(_ context.Context, id string) (*domainteam.Team, error) {
			return &domainteam.Team{ID: id, Name: "Test Team"}, nil
		},
	}
	svc := NewService(teamRepo, memberRepoWithMember(domainteam.RoleMember))

	tm, err := svc.GetTeam(context.Background(), "user-1", "team-1")
	assert.NoError(t, err)
	assert.Equal(t, "team-1", tm.ID)
}

// IDOR guard: a non-member must not be able to read a team by ID.
func TestService_GetTeam_NonMemberForbidden(t *testing.T) {
	teamRepo := &mockTeamRepo{
		findByIDFn: func(_ context.Context, id string) (*domainteam.Team, error) {
			return &domainteam.Team{ID: id, Name: "Secret Team"}, nil
		},
	}
	// default mockMemberRepo: FindByTeamAndUser returns ErrNotFound (not a member)
	svc := NewService(teamRepo, &mockMemberRepo{})

	_, err := svc.GetTeam(context.Background(), "outsider", "team-1")
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestService_GetTeam_NotFound(t *testing.T) {
	// Member exists, but the team itself is missing.
	svc := NewService(&mockTeamRepo{}, memberRepoWithMember(domainteam.RoleMember))

	_, err := svc.GetTeam(context.Background(), "user-1", "nonexistent")
	assert.Error(t, err)
}

// --- GetTeamBySlug ---

func TestService_GetTeamBySlug_Success(t *testing.T) {
	teamRepo := &mockTeamRepo{
		findBySlugFn: func(_ context.Context, slug string) (*domainteam.Team, error) {
			return &domainteam.Team{ID: "team-1", Slug: slug}, nil
		},
	}
	svc := NewService(teamRepo, memberRepoWithMember(domainteam.RoleMember))

	tm, err := svc.GetTeamBySlug(context.Background(), "user-1", "my-team")
	assert.NoError(t, err)
	assert.Equal(t, "my-team", tm.Slug)
}

// IDOR guard: a non-member must not be able to read a team by slug.
func TestService_GetTeamBySlug_NonMemberForbidden(t *testing.T) {
	teamRepo := &mockTeamRepo{
		findBySlugFn: func(_ context.Context, slug string) (*domainteam.Team, error) {
			return &domainteam.Team{ID: "team-1", Slug: slug}, nil
		},
	}
	svc := NewService(teamRepo, &mockMemberRepo{})

	_, err := svc.GetTeamBySlug(context.Background(), "outsider", "my-team")
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

// --- UpdateTeam ---

func TestService_UpdateTeam_Success(t *testing.T) {
	teamRepo := &mockTeamRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainteam.Team, error) {
			return &domainteam.Team{ID: "team-1", Name: "Old Name"}, nil
		},
	}
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleOwner}, nil
		},
	}
	svc := NewService(teamRepo, memberRepo)

	tm, err := svc.UpdateTeam(context.Background(), "user-1", "team-1", "New Name")
	assert.NoError(t, err)
	assert.Equal(t, "New Name", tm.Name)
	assert.Equal(t, "new-name", tm.Slug)
}

func TestService_UpdateTeam_AdminCanUpdate(t *testing.T) {
	teamRepo := &mockTeamRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainteam.Team, error) {
			return &domainteam.Team{ID: "team-1", Name: "Old Name"}, nil
		},
	}
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleAdmin}, nil
		},
	}
	svc := NewService(teamRepo, memberRepo)

	tm, err := svc.UpdateTeam(context.Background(), "user-1", "team-1", "New Name")
	assert.NoError(t, err)
	assert.Equal(t, "New Name", tm.Name)
}

func TestService_UpdateTeam_MemberForbidden(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleMember}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	_, err := svc.UpdateTeam(context.Background(), "user-1", "team-1", "New Name")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestService_UpdateTeam_NotAMember(t *testing.T) {
	svc := NewService(&mockTeamRepo{}, &mockMemberRepo{})

	_, err := svc.UpdateTeam(context.Background(), "user-1", "team-1", "New Name")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

// --- DeleteTeam ---

func TestService_DeleteTeam_Success(t *testing.T) {
	var deletedID string
	teamRepo := &mockTeamRepo{
		deleteFn: func(_ context.Context, id string) error {
			deletedID = id
			return nil
		},
	}
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleOwner}, nil
		},
	}
	svc := NewService(teamRepo, memberRepo)

	err := svc.DeleteTeam(context.Background(), "user-1", "team-1")
	assert.NoError(t, err)
	assert.Equal(t, "team-1", deletedID)
}

func TestService_DeleteTeam_AdminForbidden(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleAdmin}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.DeleteTeam(context.Background(), "user-1", "team-1")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestService_DeleteTeam_MemberForbidden(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleMember}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.DeleteTeam(context.Background(), "user-1", "team-1")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

// --- ListUserTeams ---

func TestService_ListUserTeams_Success(t *testing.T) {
	teamRepo := &mockTeamRepo{
		listByUserFn: func(_ context.Context, _ string) ([]*domainteam.Team, error) {
			return []*domainteam.Team{
				{ID: "team-1", Name: "Team A"},
				{ID: "team-2", Name: "Team B"},
			}, nil
		},
	}
	svc := NewService(teamRepo, &mockMemberRepo{})

	teams, err := svc.ListUserTeams(context.Background(), "user-1")
	assert.NoError(t, err)
	assert.Len(t, teams, 2)
}

// --- InviteMember ---

func TestService_InviteMember_Success(t *testing.T) {
	teamRepo := &mockTeamRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainteam.Team, error) {
			return &domainteam.Team{ID: "team-1", MaxMembers: 5}, nil
		},
	}
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleOwner}, nil
		},
		countByTeamIDFn: func(_ context.Context, _ string) (int, error) {
			return 2, nil
		},
	}
	svc := NewService(teamRepo, memberRepo)

	invite, err := svc.InviteMember(context.Background(), "user-1", "team-1", "new@example.com", domainteam.RoleMember)
	assert.NoError(t, err)
	assert.Equal(t, "new@example.com", invite.InvitedEmail)
	assert.Equal(t, domainteam.RoleMember, invite.Role)
	assert.Equal(t, domainteam.InvitePending, invite.InviteStatus)
	assert.Len(t, invite.InviteToken, 64)
}

func TestService_InviteMember_TeamFull(t *testing.T) {
	teamRepo := &mockTeamRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainteam.Team, error) {
			return &domainteam.Team{ID: "team-1", MaxMembers: 3}, nil
		},
	}
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleOwner}, nil
		},
		countByTeamIDFn: func(_ context.Context, _ string) (int, error) {
			return 3, nil
		},
	}
	svc := NewService(teamRepo, memberRepo)

	_, err := svc.InviteMember(context.Background(), "user-1", "team-1", "new@example.com", domainteam.RoleMember)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestService_InviteMember_NoPermission(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleMember}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	_, err := svc.InviteMember(context.Background(), "user-1", "team-1", "new@example.com", domainteam.RoleMember)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestService_InviteMember_NotAMember(t *testing.T) {
	svc := NewService(&mockTeamRepo{}, &mockMemberRepo{})

	_, err := svc.InviteMember(context.Background(), "user-1", "team-1", "new@example.com", domainteam.RoleMember)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestService_InviteMember_AdminCanInvite(t *testing.T) {
	teamRepo := &mockTeamRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainteam.Team, error) {
			return &domainteam.Team{ID: "team-1", MaxMembers: 5}, nil
		},
	}
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleAdmin}, nil
		},
		countByTeamIDFn: func(_ context.Context, _ string) (int, error) {
			return 1, nil
		},
	}
	svc := NewService(teamRepo, memberRepo)

	invite, err := svc.InviteMember(context.Background(), "user-1", "team-1", "new@example.com", domainteam.RoleMember)
	assert.NoError(t, err)
	assert.NotNil(t, invite)
}

// --- AcceptInvite ---

func TestService_AcceptInvite_Success(t *testing.T) {
	var updatedMember *domainteam.TeamMember
	memberRepo := &mockMemberRepo{
		findByInviteTokenFn: func(_ context.Context, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{
				ID:           "member-1",
				TeamID:       "team-1",
				Role:         domainteam.RoleMember,
				InviteStatus: domainteam.InvitePending,
				InviteToken:  "valid-token",
			}, nil
		},
		updateFn: func(_ context.Context, member *domainteam.TeamMember) error {
			updatedMember = member
			return nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	member, err := svc.AcceptInvite(context.Background(), "user-2", "valid-token")
	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, "user-2", updatedMember.UserID)
	assert.Equal(t, domainteam.InviteAccepted, updatedMember.InviteStatus)
	assert.NotNil(t, updatedMember.JoinedAt)
}

func TestService_AcceptInvite_InvalidToken(t *testing.T) {
	svc := NewService(&mockTeamRepo{}, &mockMemberRepo{})

	_, err := svc.AcceptInvite(context.Background(), "user-2", "invalid-token")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidToken)
}

func TestService_AcceptInvite_AlreadyAccepted(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByInviteTokenFn: func(_ context.Context, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{
				InviteStatus: domainteam.InviteAccepted,
			}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	_, err := svc.AcceptInvite(context.Background(), "user-2", "used-token")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

// --- DeclineInvite ---

func TestService_DeclineInvite_Success(t *testing.T) {
	var updatedMember *domainteam.TeamMember
	memberRepo := &mockMemberRepo{
		findByInviteTokenFn: func(_ context.Context, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{
				ID:           "member-1",
				InviteStatus: domainteam.InvitePending,
			}, nil
		},
		updateFn: func(_ context.Context, member *domainteam.TeamMember) error {
			updatedMember = member
			return nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.DeclineInvite(context.Background(), "valid-token")
	assert.NoError(t, err)
	assert.Equal(t, domainteam.InviteDeclined, updatedMember.InviteStatus)
}

func TestService_DeclineInvite_InvalidToken(t *testing.T) {
	svc := NewService(&mockTeamRepo{}, &mockMemberRepo{})

	err := svc.DeclineInvite(context.Background(), "invalid-token")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidToken)
}

func TestService_DeclineInvite_AlreadyAccepted(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByInviteTokenFn: func(_ context.Context, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{
				InviteStatus: domainteam.InviteAccepted,
			}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.DeclineInvite(context.Background(), "used-token")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

// --- RemoveMember ---

func TestService_RemoveMember_OwnerRemovesMember(t *testing.T) {
	var removedTeamID, removedUserID string
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, teamID, userID string) (*domainteam.TeamMember, error) {
			if userID == "owner-1" {
				return &domainteam.TeamMember{Role: domainteam.RoleOwner}, nil
			}
			if userID == "member-1" {
				return &domainteam.TeamMember{Role: domainteam.RoleMember}, nil
			}
			return nil, domain.ErrNotFound
		},
		removeFn: func(_ context.Context, teamID, userID string) error {
			removedTeamID = teamID
			removedUserID = userID
			return nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.RemoveMember(context.Background(), "owner-1", "team-1", "member-1")
	assert.NoError(t, err)
	assert.Equal(t, "team-1", removedTeamID)
	assert.Equal(t, "member-1", removedUserID)
}

func TestService_RemoveMember_CannotRemoveOwner(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, userID string) (*domainteam.TeamMember, error) {
			// Both are owners (the remover and target are the owner)
			return &domainteam.TeamMember{Role: domainteam.RoleOwner}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.RemoveMember(context.Background(), "owner-1", "team-1", "owner-1")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestService_RemoveMember_MemberNoPermission(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, userID string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleMember}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.RemoveMember(context.Background(), "member-1", "team-1", "member-2")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestService_RemoveMember_AdminCannotRemoveAdmin(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, userID string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleAdmin}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.RemoveMember(context.Background(), "admin-1", "team-1", "admin-2")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestService_RemoveMember_AdminRemovesMember(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, userID string) (*domainteam.TeamMember, error) {
			if userID == "admin-1" {
				return &domainteam.TeamMember{Role: domainteam.RoleAdmin}, nil
			}
			return &domainteam.TeamMember{Role: domainteam.RoleMember}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.RemoveMember(context.Background(), "admin-1", "team-1", "member-1")
	assert.NoError(t, err)
}

// --- UpdateMemberRole ---

func TestService_UpdateMemberRole_OwnerChangesRole(t *testing.T) {
	var updatedMember *domainteam.TeamMember
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, userID string) (*domainteam.TeamMember, error) {
			if userID == "owner-1" {
				return &domainteam.TeamMember{Role: domainteam.RoleOwner}, nil
			}
			return &domainteam.TeamMember{Role: domainteam.RoleMember, UserID: "member-1"}, nil
		},
		updateFn: func(_ context.Context, member *domainteam.TeamMember) error {
			updatedMember = member
			return nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.UpdateMemberRole(context.Background(), "owner-1", "team-1", "member-1", domainteam.RoleAdmin)
	assert.NoError(t, err)
	assert.Equal(t, domainteam.RoleAdmin, updatedMember.Role)
}

func TestService_UpdateMemberRole_AdminForbidden(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, userID string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleAdmin}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.UpdateMemberRole(context.Background(), "admin-1", "team-1", "member-1", domainteam.RoleAdmin)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestService_UpdateMemberRole_MemberForbidden(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleMember}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.UpdateMemberRole(context.Background(), "member-1", "team-1", "member-2", domainteam.RoleAdmin)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestService_UpdateMemberRole_CannotChangeOwnerRole(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, userID string) (*domainteam.TeamMember, error) {
			// Both calls return owner (updater is owner, target is also owner)
			return &domainteam.TeamMember{Role: domainteam.RoleOwner}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.UpdateMemberRole(context.Background(), "owner-1", "team-1", "owner-1", domainteam.RoleAdmin)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestService_UpdateMemberRole_CannotAssignOwner(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, userID string) (*domainteam.TeamMember, error) {
			if userID == "owner-1" {
				return &domainteam.TeamMember{Role: domainteam.RoleOwner}, nil
			}
			return &domainteam.TeamMember{Role: domainteam.RoleMember}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.UpdateMemberRole(context.Background(), "owner-1", "team-1", "member-1", domainteam.RoleOwner)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestService_UpdateMemberRole_InvalidRole(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleOwner}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.UpdateMemberRole(context.Background(), "owner-1", "team-1", "member-1", domainteam.Role("superadmin"))
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

// --- ListMembers ---

func TestService_ListMembers_Success(t *testing.T) {
	members := []*domainteam.TeamMember{
		{ID: "m1", UserID: "user-1", Role: domainteam.RoleOwner},
		{ID: "m2", UserID: "user-2", Role: domainteam.RoleMember},
	}
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleMember}, nil
		},
		listByTeamIDFn: func(_ context.Context, _ string, _, _ int) ([]*domainteam.TeamMember, int, error) {
			return members, 2, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	result, total, err := svc.ListMembers(context.Background(), "user-1", "team-1", 0, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, total)
}

func TestService_ListMembers_NotAMember(t *testing.T) {
	svc := NewService(&mockTeamRepo{}, &mockMemberRepo{})

	_, _, err := svc.ListMembers(context.Background(), "outsider", "team-1", 0, 20)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

// --- LeaveTeam ---

func TestService_LeaveTeam_MemberLeaves(t *testing.T) {
	var removedTeamID, removedUserID string
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleMember}, nil
		},
		removeFn: func(_ context.Context, teamID, userID string) error {
			removedTeamID = teamID
			removedUserID = userID
			return nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.LeaveTeam(context.Background(), "user-1", "team-1")
	assert.NoError(t, err)
	assert.Equal(t, "team-1", removedTeamID)
	assert.Equal(t, "user-1", removedUserID)
}

func TestService_LeaveTeam_AdminLeaves(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleAdmin}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.LeaveTeam(context.Background(), "admin-1", "team-1")
	assert.NoError(t, err)
}

func TestService_LeaveTeam_OwnerCannotLeave(t *testing.T) {
	memberRepo := &mockMemberRepo{
		findByTeamAndUserFn: func(_ context.Context, _, _ string) (*domainteam.TeamMember, error) {
			return &domainteam.TeamMember{Role: domainteam.RoleOwner}, nil
		},
	}
	svc := NewService(&mockTeamRepo{}, memberRepo)

	err := svc.LeaveTeam(context.Background(), "owner-1", "team-1")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestService_LeaveTeam_NotAMember(t *testing.T) {
	svc := NewService(&mockTeamRepo{}, &mockMemberRepo{})

	err := svc.LeaveTeam(context.Background(), "outsider", "team-1")
	assert.Error(t, err)
}
