package admin

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

type Service struct {
	users         repository.UserRepository
	subscriptions repository.SubscriptionRepository
}

func NewService(users repository.UserRepository, subscriptions repository.SubscriptionRepository) *Service {
	return &Service{users: users, subscriptions: subscriptions}
}

// Dashboard, Analytics, ManageUsers will be implemented
