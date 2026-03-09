package user

import (
	"context"

	domainuser "github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
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

func (s *Service) DeleteAccount(ctx context.Context, userID string) error {
	return s.users.Delete(ctx, userID)
}
