package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hassad/boilerplateSaaS/backend/internal/adapter/postgres"
	"github.com/hassad/boilerplateSaaS/backend/internal/config"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/pkg/hash"
)

func main() {
	cfg := config.Load()
	db := postgres.NewDB(cfg.DatabaseURL)
	defer db.Close()

	ctx := context.Background()
	userRepo := postgres.NewUserRepository(db)

	seedAdmin(ctx, userRepo)

	fmt.Println("Seed completed.")
}

func seedAdmin(ctx context.Context, repo *postgres.UserRepository) {
	_, err := repo.FindByEmail(ctx, "admin@cleansaas.dev")
	if err == nil {
		fmt.Println("Admin user already exists, skipping.")
		return
	}

	hashed, err := hash.Password("admin123")
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	admin, err := user.New("admin@cleansaas.dev", "Admin", hashed)
	if err != nil {
		log.Fatalf("failed to create admin user: %v", err)
	}
	admin.Role = user.RoleAdmin
	admin.EmailVerified = true

	if err := repo.Create(ctx, admin); err != nil {
		log.Fatalf("failed to insert admin: %v", err)
	}

	fmt.Println("Admin user created: admin@cleansaas.dev / admin123")
}
