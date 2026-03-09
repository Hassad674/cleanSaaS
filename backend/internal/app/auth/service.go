package auth

import (
	"context"
	"errors"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
	"github.com/hassad/boilerplateSaaS/backend/pkg/hash"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
)

type Service struct {
	users    repository.UserRepository
	email    service.EmailService
	jwtMaker *jwt.Maker
}

func NewService(users repository.UserRepository, email service.EmailService, jwtMaker *jwt.Maker) *Service {
	return &Service{users: users, email: email, jwtMaker: jwtMaker}
}

func (s *Service) Register(ctx context.Context, email, name, password string) (*user.User, string, error) {
	_, err := s.users.FindByEmail(ctx, email)
	if err == nil {
		return nil, "", domain.ErrAlreadyExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, "", err
	}

	hashed, err := hash.Password(password)
	if err != nil {
		return nil, "", err
	}

	u, err := user.New(email, name, hashed)
	if err != nil {
		return nil, "", err
	}

	if err := s.users.Create(ctx, u); err != nil {
		return nil, "", err
	}

	token, err := s.jwtMaker.Generate(u.ID, string(u.Role))
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*user.User, string, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", domain.ErrUnauthorized
	}

	if !hash.Check(password, u.PasswordHash) {
		return nil, "", domain.ErrUnauthorized
	}

	token, err := s.jwtMaker.Generate(u.ID, string(u.Role))
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}

func (s *Service) OAuthCallback(ctx context.Context, oauthUser *service.OAuthUser) (*user.User, string, error) {
	u, err := s.users.FindByProvider(ctx, oauthUser.Provider, oauthUser.ProviderID)
	if errors.Is(err, domain.ErrNotFound) {
		u = user.NewOAuth(oauthUser.Email, oauthUser.Name, oauthUser.Provider, oauthUser.ProviderID)
		u.AvatarURL = oauthUser.AvatarURL
		if err := s.users.Create(ctx, u); err != nil {
			return nil, "", err
		}
	} else if err != nil {
		return nil, "", err
	}

	token, err := s.jwtMaker.Generate(u.ID, string(u.Role))
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}
