package user

import (
	"context"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainuser "github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/pkg/hash"
)

type Service struct {
	users repository.UserRepository
}

func NewService(users repository.UserRepository) *Service {
	return &Service{users: users}
}

func (s *Service) GetProfile(ctx context.Context, userID string) (*domainuser.User, error) {
	return s.users.FindByID(ctx, userID)
}

func (s *Service) UpdateProfile(ctx context.Context, userID, name, avatarURL string) (*domainuser.User, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if name != "" {
		u.Name = name
	}
	if avatarURL != "" {
		u.AvatarURL = avatarURL
	}

	if err := s.users.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("finding user: %w", err)
	}

	if !hash.Check(oldPassword, u.PasswordHash) {
		return domain.ErrUnauthorized
	}

	if oldPassword == newPassword {
		return fmt.Errorf("new password must be different: %w", domain.ErrValidation)
	}

	hashed, err := hash.Password(newPassword)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	u.PasswordHash = hashed
	if err := s.users.Update(ctx, u); err != nil {
		return fmt.Errorf("updating password: %w", err)
	}

	return nil
}

func (s *Service) DeleteAccount(ctx context.Context, userID string) error {
	return s.users.Delete(ctx, userID)
}
