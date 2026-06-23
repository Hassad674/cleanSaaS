package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/team"
)

// TeamRepository implements repository.TeamRepository using PostgreSQL.
// It holds a DBTX (not a concrete *sql.DB) so the same code runs against either the
// connection pool or an open transaction — that is what lets it join a shared tx.
type TeamRepository struct {
	db DBTX
}

// NewTeamRepository creates a new PostgreSQL-backed team repository bound to the pool.
func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

// newTeamRepositoryTx creates a team repository bound to an open transaction, so its
// writes participate in that transaction. Used by TxManager.
func newTeamRepositoryTx(tx DBTX) *TeamRepository {
	return &TeamRepository{db: tx}
}

func (r *TeamRepository) Create(ctx context.Context, t *team.Team) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO teams (name, slug, owner_id, avatar_url, plan, max_members, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		t.Name, t.Slug, t.OwnerID, t.AvatarURL, t.Plan, t.MaxMembers, t.CreatedAt, t.UpdatedAt,
	).Scan(&t.ID)
	if err != nil {
		return fmt.Errorf("creating team: %w", err)
	}
	return nil
}

func (r *TeamRepository) FindByID(ctx context.Context, id string) (*team.Team, error) {
	t := &team.Team{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, slug, owner_id, avatar_url, plan, max_members, created_at, updated_at
		 FROM teams WHERE id = $1`, id,
	).Scan(&t.ID, &t.Name, &t.Slug, &t.OwnerID, &t.AvatarURL, &t.Plan, &t.MaxMembers, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding team by ID: %w", err)
	}
	return t, nil
}

func (r *TeamRepository) FindBySlug(ctx context.Context, slug string) (*team.Team, error) {
	t := &team.Team{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, slug, owner_id, avatar_url, plan, max_members, created_at, updated_at
		 FROM teams WHERE slug = $1`, slug,
	).Scan(&t.ID, &t.Name, &t.Slug, &t.OwnerID, &t.AvatarURL, &t.Plan, &t.MaxMembers, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding team by slug: %w", err)
	}
	return t, nil
}

func (r *TeamRepository) Update(ctx context.Context, t *team.Team) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE teams SET name = $1, slug = $2, avatar_url = $3, plan = $4, max_members = $5, updated_at = $6
		 WHERE id = $7`,
		t.Name, t.Slug, t.AvatarURL, t.Plan, t.MaxMembers, t.UpdatedAt, t.ID,
	)
	if err != nil {
		return fmt.Errorf("updating team: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TeamRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM teams WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting team: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TeamRepository) ListByUserID(ctx context.Context, userID string) ([]*team.Team, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT t.id, t.name, t.slug, t.owner_id, t.avatar_url, t.plan, t.max_members, t.created_at, t.updated_at
		 FROM teams t
		 JOIN team_members tm ON t.id = tm.team_id
		 WHERE tm.user_id = $1 AND tm.invite_status = 'accepted'
		 ORDER BY t.created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing teams by user: %w", err)
	}
	defer rows.Close()

	var teams []*team.Team
	for rows.Next() {
		t := &team.Team{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.OwnerID, &t.AvatarURL, &t.Plan, &t.MaxMembers, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning team: %w", err)
		}
		teams = append(teams, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating teams: %w", err)
	}
	return teams, nil
}

// TeamMemberRepository implements repository.TeamMemberRepository using PostgreSQL.
// It holds a DBTX so its writes can run on either the pool or an open transaction.
type TeamMemberRepository struct {
	db DBTX
}

// NewTeamMemberRepository creates a new PostgreSQL-backed team member repository bound to the pool.
func NewTeamMemberRepository(db *sql.DB) *TeamMemberRepository {
	return &TeamMemberRepository{db: db}
}

// newTeamMemberRepositoryTx creates a team member repository bound to an open
// transaction, so its writes participate in that transaction. Used by TxManager.
func newTeamMemberRepositoryTx(tx DBTX) *TeamMemberRepository {
	return &TeamMemberRepository{db: tx}
}

func (r *TeamMemberRepository) Add(ctx context.Context, member *team.TeamMember) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO team_members (team_id, user_id, role, invited_email, invite_token, invite_status, joined_at, created_at)
		 VALUES ($1, NULLIF($2, '')::uuid, $3, NULLIF($4, ''), NULLIF($5, ''), $6, $7, $8)
		 RETURNING id`,
		member.TeamID, member.UserID, string(member.Role), member.InvitedEmail, member.InviteToken, string(member.InviteStatus), member.JoinedAt, member.CreatedAt,
	).Scan(&member.ID)
	if err != nil {
		return fmt.Errorf("adding team member: %w", err)
	}
	return nil
}

func (r *TeamMemberRepository) FindByID(ctx context.Context, id string) (*team.TeamMember, error) {
	m := &team.TeamMember{}
	var userID sql.NullString
	var invitedEmail sql.NullString
	var inviteToken sql.NullString
	var joinedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, team_id, user_id, role, invited_email, invite_token, invite_status, joined_at, created_at
		 FROM team_members WHERE id = $1`, id,
	).Scan(&m.ID, &m.TeamID, &userID, &m.Role, &invitedEmail, &inviteToken, &m.InviteStatus, &joinedAt, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding team member by ID: %w", err)
	}

	if userID.Valid {
		m.UserID = userID.String
	}
	if invitedEmail.Valid {
		m.InvitedEmail = invitedEmail.String
	}
	if inviteToken.Valid {
		m.InviteToken = inviteToken.String
	}
	if joinedAt.Valid {
		m.JoinedAt = &joinedAt.Time
	}
	return m, nil
}

func (r *TeamMemberRepository) FindByTeamAndUser(ctx context.Context, teamID, userID string) (*team.TeamMember, error) {
	m := &team.TeamMember{}
	var nullUserID sql.NullString
	var invitedEmail sql.NullString
	var inviteToken sql.NullString
	var joinedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, team_id, user_id, role, invited_email, invite_token, invite_status, joined_at, created_at
		 FROM team_members WHERE team_id = $1 AND user_id = $2`, teamID, userID,
	).Scan(&m.ID, &m.TeamID, &nullUserID, &m.Role, &invitedEmail, &inviteToken, &m.InviteStatus, &joinedAt, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding team member by team and user: %w", err)
	}

	if nullUserID.Valid {
		m.UserID = nullUserID.String
	}
	if invitedEmail.Valid {
		m.InvitedEmail = invitedEmail.String
	}
	if inviteToken.Valid {
		m.InviteToken = inviteToken.String
	}
	if joinedAt.Valid {
		m.JoinedAt = &joinedAt.Time
	}
	return m, nil
}

func (r *TeamMemberRepository) FindByInviteToken(ctx context.Context, token string) (*team.TeamMember, error) {
	m := &team.TeamMember{}
	var userID sql.NullString
	var invitedEmail sql.NullString
	var inviteToken sql.NullString
	var joinedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, team_id, user_id, role, invited_email, invite_token, invite_status, joined_at, created_at
		 FROM team_members WHERE invite_token = $1`, token,
	).Scan(&m.ID, &m.TeamID, &userID, &m.Role, &invitedEmail, &inviteToken, &m.InviteStatus, &joinedAt, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding team member by invite token: %w", err)
	}

	if userID.Valid {
		m.UserID = userID.String
	}
	if invitedEmail.Valid {
		m.InvitedEmail = invitedEmail.String
	}
	if inviteToken.Valid {
		m.InviteToken = inviteToken.String
	}
	if joinedAt.Valid {
		m.JoinedAt = &joinedAt.Time
	}
	return m, nil
}

func (r *TeamMemberRepository) Update(ctx context.Context, member *team.TeamMember) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE team_members SET user_id = NULLIF($1, '')::uuid, role = $2, invite_status = $3, joined_at = $4
		 WHERE id = $5`,
		member.UserID, string(member.Role), string(member.InviteStatus), member.JoinedAt, member.ID,
	)
	if err != nil {
		return fmt.Errorf("updating team member: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TeamMemberRepository) Remove(ctx context.Context, teamID, userID string) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM team_members WHERE team_id = $1 AND user_id = $2`, teamID, userID,
	)
	if err != nil {
		return fmt.Errorf("removing team member: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TeamMemberRepository) ListByTeamID(ctx context.Context, teamID string, offset, limit int) ([]*team.TeamMember, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM team_members WHERE team_id = $1`, teamID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting team members: %w", err)
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, team_id, user_id, role, invited_email, invite_token, invite_status, joined_at, created_at
		 FROM team_members WHERE team_id = $1
		 ORDER BY created_at ASC LIMIT $2 OFFSET $3`,
		teamID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("listing team members: %w", err)
	}
	defer rows.Close()

	var members []*team.TeamMember
	for rows.Next() {
		m := &team.TeamMember{}
		var userID sql.NullString
		var invitedEmail sql.NullString
		var inviteToken sql.NullString
		var joinedAt sql.NullTime

		if err := rows.Scan(&m.ID, &m.TeamID, &userID, &m.Role, &invitedEmail, &inviteToken, &m.InviteStatus, &joinedAt, &m.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning team member: %w", err)
		}

		if userID.Valid {
			m.UserID = userID.String
		}
		if invitedEmail.Valid {
			m.InvitedEmail = invitedEmail.String
		}
		if inviteToken.Valid {
			m.InviteToken = inviteToken.String
		}
		if joinedAt.Valid {
			m.JoinedAt = &joinedAt.Time
		}

		members = append(members, m)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating team members: %w", err)
	}
	return members, total, nil
}

func (r *TeamMemberRepository) CountByTeamID(ctx context.Context, teamID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM team_members WHERE team_id = $1`, teamID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting team members: %w", err)
	}
	return count, nil
}
