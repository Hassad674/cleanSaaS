// Package org provides organization use cases: resolving a caller's active
// tenant and looking organizations up. It is the application seam the auth
// middleware uses (via an injected resolver) to turn an authenticated user into an
// authorized active organization for the request.
package org

import (
	"context"
	"errors"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainorg "github.com/hassad/boilerplateSaaS/backend/internal/domain/org"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

// Service orchestrates organization use cases.
type Service struct {
	orgs    repository.OrganizationRepository
	members repository.OrganizationMemberRepository
}

// NewService wires the organization service.
func NewService(orgs repository.OrganizationRepository, members repository.OrganizationMemberRepository) *Service {
	return &Service{orgs: orgs, members: members}
}

// ResolveActiveOrg returns the organization id the request should be scoped to.
//
//   - If the token carries an explicit org, the user MUST be a member of it; an
//     unauthorized org is rejected (a user cannot scope a request to a tenant they
//     don't belong to). This closes the "forge the org claim" hole.
//   - Otherwise the user's default/personal organization is used.
//   - If neither resolves, it returns "" so the tenant request path fails closed.
func (s *Service) ResolveActiveOrg(ctx context.Context, userID, tokenOrgID string) (string, error) {
	if userID == "" {
		return "", nil
	}

	if tokenOrgID != "" {
		ok, err := s.members.IsMember(ctx, tokenOrgID, userID)
		if err != nil {
			return "", fmt.Errorf("checking org membership: %w", err)
		}
		if !ok {
			return "", domain.ErrForbidden
		}
		return tokenOrgID, nil
	}

	o, err := s.orgs.FindDefaultForUser(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", nil
		}
		return "", fmt.Errorf("resolving default org: %w", err)
	}
	return o.ID, nil
}

// GetByID returns an organization by id.
func (s *Service) GetByID(ctx context.Context, id string) (*domainorg.Organization, error) {
	return s.orgs.FindByID(ctx, id)
}
