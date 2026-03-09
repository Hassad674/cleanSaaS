package service

import "context"

type OAuthUser struct {
	Email      string
	Name       string
	AvatarURL  string
	Provider   string
	ProviderID string
}

type OAuthProvider interface {
	GetAuthURL(state string) string
	Exchange(ctx context.Context, code string) (*OAuthUser, error)
}
