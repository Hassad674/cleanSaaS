package repository

import (
	"context"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
)

type UserRepository interface {
	Create(ctx context.Context, u *user.User) error
	FindByID(ctx context.Context, id string) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindByProvider(ctx context.Context, provider, providerID string) (*user.User, error)
	FindByStripeID(ctx context.Context, stripeID string) (*user.User, error)
	Update(ctx context.Context, u *user.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*user.User, int, error)
	Search(ctx context.Context, query string, offset, limit int) ([]*user.User, int, error)
	Count(ctx context.Context) (int, error)
}
