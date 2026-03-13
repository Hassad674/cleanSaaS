package repository

import (
	"context"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/team"
)

// TeamRepository defines persistence operations for teams.
type TeamRepository interface {
	Create(ctx context.Context, t *team.Team) error
	FindByID(ctx context.Context, id string) (*team.Team, error)
	FindBySlug(ctx context.Context, slug string) (*team.Team, error)
	Update(ctx context.Context, t *team.Team) error
	Delete(ctx context.Context, id string) error
	ListByUserID(ctx context.Context, userID string) ([]*team.Team, error)
}

// TeamMemberRepository defines persistence operations for team members.
type TeamMemberRepository interface {
	Add(ctx context.Context, member *team.TeamMember) error
	FindByID(ctx context.Context, id string) (*team.TeamMember, error)
	FindByTeamAndUser(ctx context.Context, teamID, userID string) (*team.TeamMember, error)
	FindByInviteToken(ctx context.Context, token string) (*team.TeamMember, error)
	Update(ctx context.Context, member *team.TeamMember) error
	Remove(ctx context.Context, teamID, userID string) error
	ListByTeamID(ctx context.Context, teamID string, offset, limit int) ([]*team.TeamMember, int, error)
	CountByTeamID(ctx context.Context, teamID string) (int, error)
}
