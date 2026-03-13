package response

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/team"
)

type TeamResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	OwnerID    string `json:"owner_id"`
	AvatarURL  string `json:"avatar_url"`
	Plan       string `json:"plan"`
	MaxMembers int    `json:"max_members"`
	CreatedAt  string `json:"created_at"`
}

type TeamMemberResponse struct {
	ID           string  `json:"id"`
	TeamID       string  `json:"team_id"`
	UserID       string  `json:"user_id"`
	Role         string  `json:"role"`
	InvitedEmail string  `json:"invited_email,omitempty"`
	InviteStatus string  `json:"invite_status"`
	JoinedAt     *string `json:"joined_at,omitempty"`
	CreatedAt    string  `json:"created_at"`
}

func TeamFromDomain(t *team.Team) TeamResponse {
	return TeamResponse{
		ID:         t.ID,
		Name:       t.Name,
		Slug:       t.Slug,
		OwnerID:    t.OwnerID,
		AvatarURL:  t.AvatarURL,
		Plan:       t.Plan,
		MaxMembers: t.MaxMembers,
		CreatedAt:  t.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func TeamMemberFromDomain(m *team.TeamMember) TeamMemberResponse {
	resp := TeamMemberResponse{
		ID:           m.ID,
		TeamID:       m.TeamID,
		UserID:       m.UserID,
		Role:         string(m.Role),
		InvitedEmail: m.InvitedEmail,
		InviteStatus: string(m.InviteStatus),
		CreatedAt:    m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if m.JoinedAt != nil {
		s := m.JoinedAt.Format("2006-01-02T15:04:05Z")
		resp.JoinedAt = &s
	}
	return resp
}
