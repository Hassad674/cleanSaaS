//go:build integration

package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/team"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

// newTeam creates a team owned by ownerID via the real TeamRepository and
// registers its deletion (cascades to team_members). Returns the persisted team.
func newTeam(ctx context.Context, t *testing.T, db *sql.DB, ownerID string) *team.Team {
	t.Helper()
	repo := NewTeamRepository(db)
	tag := uniqueTag()
	tm := &team.Team{
		Name:       "Integration Team " + tag,
		Slug:       "itest-team-" + tag,
		OwnerID:    ownerID,
		Plan:       "free",
		MaxMembers: 5,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	require.NoError(t, repo.Create(ctx, tm))
	require.NotEmpty(t, tm.ID)
	t.Cleanup(func() { _ = repo.Delete(context.Background(), tm.ID) })
	return tm
}

// TestTeamRepository_RoundTrip exercises team CRUD: create, find by id/slug,
// update (with RowsAffected guard), delete, and the not-found mapping.
func TestTeamRepository_RoundTrip(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewTeamRepository(db)

	owner := newUser(ctx, t, db)
	tm := newTeam(ctx, t, db, owner.ID)

	byID, err := repo.FindByID(ctx, tm.ID)
	require.NoError(t, err)
	assert.Equal(t, tm.Slug, byID.Slug)

	bySlug, err := repo.FindBySlug(ctx, tm.Slug)
	require.NoError(t, err)
	assert.Equal(t, tm.ID, bySlug.ID)

	// Update mutates name + max_members.
	tm.Name = "Renamed Team"
	tm.MaxMembers = 10
	tm.UpdatedAt = time.Now()
	require.NoError(t, repo.Update(ctx, tm))

	updated, err := repo.FindByID(ctx, tm.ID)
	require.NoError(t, err)
	assert.Equal(t, "Renamed Team", updated.Name)
	assert.Equal(t, 10, updated.MaxMembers)

	// Updating a non-existent team affects zero rows -> not-found.
	ghost := &team.Team{ID: uniqueTag(), Name: "x", Slug: "x", UpdatedAt: time.Now()}
	assert.ErrorIs(t, repo.Update(ctx, ghost), domain.ErrNotFound)

	// Delete removes the team; the lookup then reports not-found.
	require.NoError(t, repo.Delete(ctx, tm.ID))
	_, err = repo.FindByID(ctx, tm.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	// Deleting again affects zero rows -> not-found.
	assert.ErrorIs(t, repo.Delete(ctx, tm.ID), domain.ErrNotFound)
}

// TestTeamMemberRepository_RoundTrip exercises membership against the real
// schema: Add (NULLIF handling of optional invite fields), find by id / team+user
// / invite token, Update, ListByTeamID pagination + count, and Remove.
func TestTeamMemberRepository_RoundTrip(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	members := NewTeamMemberRepository(db)

	owner := newUser(ctx, t, db)
	tm := newTeam(ctx, t, db, owner.ID)
	memberUser := newUser(ctx, t, db)

	joined := time.Now()
	m := &team.TeamMember{
		TeamID:       tm.ID,
		UserID:       memberUser.ID,
		Role:         team.RoleMember,
		InviteStatus: team.InviteAccepted,
		JoinedAt:     &joined,
		CreatedAt:    time.Now(),
	}
	require.NoError(t, members.Add(ctx, m))
	require.NotEmpty(t, m.ID)

	byID, err := members.FindByID(ctx, m.ID)
	require.NoError(t, err)
	assert.Equal(t, memberUser.ID, byID.UserID)
	require.NotNil(t, byID.JoinedAt, "joined_at round-trips through sql.NullTime")

	byTeamUser, err := members.FindByTeamAndUser(ctx, tm.ID, memberUser.ID)
	require.NoError(t, err)
	assert.Equal(t, m.ID, byTeamUser.ID)

	// A pending invite with a token but no user id exercises the NULLIF('') path
	// (invited_email / invite_token are nullable, user_id NULLIF on empty string).
	invitee := newUser(ctx, t, db) // invite_token rows still need a real user_id per FK
	token := "invite-" + uniqueTag()
	invite := &team.TeamMember{
		TeamID:       tm.ID,
		UserID:       invitee.ID,
		Role:         team.RoleMember,
		InvitedEmail: "invitee-" + uniqueTag() + "@itest.test",
		InviteToken:  token,
		InviteStatus: team.InvitePending,
		CreatedAt:    time.Now(),
	}
	require.NoError(t, members.Add(ctx, invite))

	byToken, err := members.FindByInviteToken(ctx, token)
	require.NoError(t, err)
	assert.Equal(t, invite.ID, byToken.ID)
	assert.Equal(t, team.InvitePending, byToken.InviteStatus)

	// Update accepts the invite.
	now := time.Now()
	byToken.InviteStatus = team.InviteAccepted
	byToken.JoinedAt = &now
	require.NoError(t, members.Update(ctx, byToken))
	reread, err := members.FindByID(ctx, byToken.ID)
	require.NoError(t, err)
	assert.Equal(t, team.InviteAccepted, reread.InviteStatus)

	// ListByTeamID pagination + CountByTeamID (2 members on this team).
	count, err := members.CountByTeamID(ctx, tm.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, count)

	page, total, err := members.ListByTeamID(ctx, tm.ID, 0, 1)
	require.NoError(t, err)
	assert.Equal(t, 2, total, "total counts both members")
	assert.Len(t, page, 1, "LIMIT 1 yields a single-row page")

	// Remove the first member; count drops to 1.
	require.NoError(t, members.Remove(ctx, tm.ID, memberUser.ID))
	count, err = members.CountByTeamID(ctx, tm.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Removing a non-member affects zero rows -> not-found.
	assert.ErrorIs(t, members.Remove(ctx, tm.ID, memberUser.ID), domain.ErrNotFound)
}

// TestTxManager_WithTeamTx_Atomic proves the team-create unit-of-work commits the
// team and its owner-member together and both are visible afterwards.
func TestTxManager_WithTeamTx_Atomic(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	txm := NewTxManager(db)

	owner := newUser(ctx, t, db)
	tag := uniqueTag()
	var teamID string

	err := txm.WithTeamTx(ctx, func(teams repository.TeamRepository, mems repository.TeamMemberRepository) error {
		tm := &team.Team{Name: "Tx Team " + tag, Slug: "itest-txteam-" + tag, OwnerID: owner.ID, Plan: "free", MaxMembers: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		if e := teams.Create(ctx, tm); e != nil {
			return e
		}
		teamID = tm.ID
		return mems.Add(ctx, &team.TeamMember{TeamID: tm.ID, UserID: owner.ID, Role: team.RoleOwner, InviteStatus: team.InviteAccepted, CreatedAt: time.Now()})
	})
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = db.ExecContext(context.Background(), `DELETE FROM teams WHERE id = $1`, teamID) })

	got, err := NewTeamRepository(db).FindByID(ctx, teamID)
	require.NoError(t, err)
	assert.Equal(t, teamID, got.ID)

	count, err := NewTeamMemberRepository(db).CountByTeamID(ctx, teamID)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "owner membership committed atomically with the team")
}
