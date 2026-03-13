package team

import (
	"context"
	"fmt"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainteam "github.com/hassad/boilerplateSaaS/backend/internal/domain/team"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

// Service orchestrates team use cases.
type Service struct {
	teams   repository.TeamRepository
	members repository.TeamMemberRepository
}

// NewService creates a new team service with the given repositories.
func NewService(teams repository.TeamRepository, members repository.TeamMemberRepository) *Service {
	return &Service{
		teams:   teams,
		members: members,
	}
}

// CreateTeam creates a new team and adds the creating user as the owner member.
func (s *Service) CreateTeam(ctx context.Context, userID, name string) (*domainteam.Team, error) {
	t, err := domainteam.NewTeam(name, userID)
	if err != nil {
		return nil, fmt.Errorf("creating team: %w", err)
	}

	if err := s.teams.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("persisting team: %w", err)
	}

	// Add the creator as the owner member
	now := time.Now()
	ownerMember := &domainteam.TeamMember{
		TeamID:       t.ID,
		UserID:       userID,
		Role:         domainteam.RoleOwner,
		InviteStatus: domainteam.InviteAccepted,
		JoinedAt:     &now,
		CreatedAt:    now,
	}

	if err := s.members.Add(ctx, ownerMember); err != nil {
		return nil, fmt.Errorf("adding owner as member: %w", err)
	}

	return t, nil
}

// GetTeam returns a team by its ID.
func (s *Service) GetTeam(ctx context.Context, teamID string) (*domainteam.Team, error) {
	t, err := s.teams.FindByID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("finding team: %w", err)
	}
	return t, nil
}

// GetTeamBySlug returns a team by its slug.
func (s *Service) GetTeamBySlug(ctx context.Context, slug string) (*domainteam.Team, error) {
	t, err := s.teams.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("finding team by slug: %w", err)
	}
	return t, nil
}

// UpdateTeam updates the team name and slug. Only owners and admins can update.
func (s *Service) UpdateTeam(ctx context.Context, userID, teamID, name string) (*domainteam.Team, error) {
	member, err := s.members.FindByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return nil, fmt.Errorf("checking membership: %w", domain.ErrForbidden)
	}

	if !domainteam.CanManageMembers(member.Role) {
		return nil, fmt.Errorf("insufficient permissions to update team: %w", domain.ErrForbidden)
	}

	t, err := s.teams.FindByID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("finding team: %w", err)
	}

	// Validate the new name
	trimmedName := name
	if len(trimmedName) < 2 || len(trimmedName) > 50 {
		return nil, fmt.Errorf("team name must be between 2 and 50 characters: %w", domain.ErrValidation)
	}

	t.Name = name
	t.Slug = domainteam.GenerateSlug(name)
	t.UpdatedAt = time.Now()

	if err := s.teams.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("updating team: %w", err)
	}

	return t, nil
}

// DeleteTeam deletes a team. Only the owner can delete.
func (s *Service) DeleteTeam(ctx context.Context, userID, teamID string) error {
	member, err := s.members.FindByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return fmt.Errorf("checking membership: %w", domain.ErrForbidden)
	}

	if !domainteam.CanDeleteTeam(member.Role) {
		return fmt.Errorf("only the team owner can delete the team: %w", domain.ErrForbidden)
	}

	if err := s.teams.Delete(ctx, teamID); err != nil {
		return fmt.Errorf("deleting team: %w", err)
	}

	return nil
}

// ListUserTeams returns all teams a user belongs to.
func (s *Service) ListUserTeams(ctx context.Context, userID string) ([]*domainteam.Team, error) {
	teams, err := s.teams.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("listing user teams: %w", err)
	}
	return teams, nil
}

// InviteMember creates a pending invitation for a new team member.
// The inviter must be an owner or admin, and the team must not exceed max_members.
func (s *Service) InviteMember(ctx context.Context, inviterUserID, teamID, email string, role domainteam.Role) (*domainteam.TeamMember, error) {
	// Check inviter permissions
	inviter, err := s.members.FindByTeamAndUser(ctx, teamID, inviterUserID)
	if err != nil {
		return nil, fmt.Errorf("checking inviter membership: %w", domain.ErrForbidden)
	}

	if !domainteam.CanManageMembers(inviter.Role) {
		return nil, fmt.Errorf("insufficient permissions to invite members: %w", domain.ErrForbidden)
	}

	// Check max members limit
	t, err := s.teams.FindByID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("finding team: %w", err)
	}

	count, err := s.members.CountByTeamID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("counting team members: %w", err)
	}

	if count >= t.MaxMembers {
		return nil, fmt.Errorf("team has reached its maximum of %d members: %w", t.MaxMembers, domain.ErrValidation)
	}

	// Create the invite
	invite, err := domainteam.NewInvite(teamID, email, role)
	if err != nil {
		return nil, fmt.Errorf("creating invite: %w", err)
	}

	if err := s.members.Add(ctx, invite); err != nil {
		return nil, fmt.Errorf("persisting invite: %w", err)
	}

	return invite, nil
}

// AcceptInvite accepts a pending invitation using the invite token.
func (s *Service) AcceptInvite(ctx context.Context, userID, token string) (*domainteam.TeamMember, error) {
	member, err := s.members.FindByInviteToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("finding invite: %w", domain.ErrInvalidToken)
	}

	if member.InviteStatus != domainteam.InvitePending {
		return nil, fmt.Errorf("invite is no longer pending: %w", domain.ErrValidation)
	}

	now := time.Now()
	member.UserID = userID
	member.InviteStatus = domainteam.InviteAccepted
	member.JoinedAt = &now

	if err := s.members.Update(ctx, member); err != nil {
		return nil, fmt.Errorf("accepting invite: %w", err)
	}

	return member, nil
}

// DeclineInvite marks an invitation as declined.
func (s *Service) DeclineInvite(ctx context.Context, token string) error {
	member, err := s.members.FindByInviteToken(ctx, token)
	if err != nil {
		return fmt.Errorf("finding invite: %w", domain.ErrInvalidToken)
	}

	if member.InviteStatus != domainteam.InvitePending {
		return fmt.Errorf("invite is no longer pending: %w", domain.ErrValidation)
	}

	member.InviteStatus = domainteam.InviteDeclined

	if err := s.members.Update(ctx, member); err != nil {
		return fmt.Errorf("declining invite: %w", err)
	}

	return nil
}

// RemoveMember removes a member from a team.
// The owner cannot be removed. Admins cannot remove other admins.
func (s *Service) RemoveMember(ctx context.Context, removerUserID, teamID, targetUserID string) error {
	// Check remover permissions
	remover, err := s.members.FindByTeamAndUser(ctx, teamID, removerUserID)
	if err != nil {
		return fmt.Errorf("checking remover membership: %w", domain.ErrForbidden)
	}

	if !domainteam.CanManageMembers(remover.Role) {
		return fmt.Errorf("insufficient permissions to remove members: %w", domain.ErrForbidden)
	}

	// Find target member
	target, err := s.members.FindByTeamAndUser(ctx, teamID, targetUserID)
	if err != nil {
		return fmt.Errorf("finding target member: %w", err)
	}

	// Cannot remove the owner
	if target.Role == domainteam.RoleOwner {
		return fmt.Errorf("cannot remove the team owner: %w", domain.ErrForbidden)
	}

	// Admin cannot remove another admin
	if remover.Role == domainteam.RoleAdmin && target.Role == domainteam.RoleAdmin {
		return fmt.Errorf("admins cannot remove other admins: %w", domain.ErrForbidden)
	}

	if err := s.members.Remove(ctx, teamID, targetUserID); err != nil {
		return fmt.Errorf("removing member: %w", err)
	}

	return nil
}

// UpdateMemberRole changes a team member's role. Only the owner can change roles.
func (s *Service) UpdateMemberRole(ctx context.Context, updaterUserID, teamID, targetUserID string, newRole domainteam.Role) error {
	// Check updater permissions — only owner can change roles
	updater, err := s.members.FindByTeamAndUser(ctx, teamID, updaterUserID)
	if err != nil {
		return fmt.Errorf("checking updater membership: %w", domain.ErrForbidden)
	}

	if updater.Role != domainteam.RoleOwner {
		return fmt.Errorf("only the team owner can change roles: %w", domain.ErrForbidden)
	}

	if !domainteam.ValidateRole(string(newRole)) {
		return fmt.Errorf("invalid role %q: %w", newRole, domain.ErrValidation)
	}

	if newRole == domainteam.RoleOwner {
		return fmt.Errorf("cannot assign owner role through role update: %w", domain.ErrValidation)
	}

	// Find target member
	target, err := s.members.FindByTeamAndUser(ctx, teamID, targetUserID)
	if err != nil {
		return fmt.Errorf("finding target member: %w", err)
	}

	// Cannot change the owner's role
	if target.Role == domainteam.RoleOwner {
		return fmt.Errorf("cannot change the owner's role: %w", domain.ErrForbidden)
	}

	target.Role = newRole

	if err := s.members.Update(ctx, target); err != nil {
		return fmt.Errorf("updating member role: %w", err)
	}

	return nil
}

// ListMembers returns a paginated list of team members. The requester must be a member.
func (s *Service) ListMembers(ctx context.Context, userID, teamID string, offset, limit int) ([]*domainteam.TeamMember, int, error) {
	// Verify the requester is a member
	_, err := s.members.FindByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("checking membership: %w", domain.ErrForbidden)
	}

	members, total, err := s.members.ListByTeamID(ctx, teamID, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("listing members: %w", err)
	}

	return members, total, nil
}

// LeaveTeam allows a member to leave a team. The owner cannot leave (must transfer or delete).
func (s *Service) LeaveTeam(ctx context.Context, userID, teamID string) error {
	member, err := s.members.FindByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return fmt.Errorf("checking membership: %w", err)
	}

	if member.Role == domainteam.RoleOwner {
		return fmt.Errorf("owner cannot leave the team, transfer ownership or delete instead: %w", domain.ErrForbidden)
	}

	if err := s.members.Remove(ctx, teamID, userID); err != nil {
		return fmt.Errorf("leaving team: %w", err)
	}

	return nil
}
