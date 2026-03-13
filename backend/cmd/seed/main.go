package main

import (
	"context"
	"database/sql"
	"encoding/json"
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
	seedPlans(ctx, db)

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

type seedPlan struct {
	Name          string
	StripePriceID string
	PriceCents    int
	Interval      string
	Features      []string
	SortOrder     int
}

func seedPlans(ctx context.Context, db *sql.DB) {
	plans := []seedPlan{
		{
			Name:          "Free",
			StripePriceID: "price_free_placeholder",
			PriceCents:    0,
			Interval:      "month",
			Features:      []string{"1 project", "Basic analytics", "Community support"},
			SortOrder:     0,
		},
		{
			Name:          "Pro",
			StripePriceID: "price_pro_placeholder",
			PriceCents:    1900,
			Interval:      "month",
			Features:      []string{"Unlimited projects", "Advanced analytics", "AI chat", "File storage (10GB)", "Priority support"},
			SortOrder:     1,
		},
		{
			Name:          "Enterprise",
			StripePriceID: "price_enterprise_placeholder",
			PriceCents:    4900,
			Interval:      "month",
			Features:      []string{"Everything in Pro", "Unlimited storage", "Custom integrations", "Dedicated support", "SLA guarantee"},
			SortOrder:     2,
		},
	}

	for _, p := range plans {
		var exists bool
		err := db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM plans WHERE name = $1)`, p.Name).Scan(&exists)
		if err != nil {
			log.Fatalf("checking plan existence: %v", err)
		}
		if exists {
			fmt.Printf("Plan '%s' already exists, skipping.\n", p.Name)
			continue
		}

		featuresJSON, _ := json.Marshal(p.Features)
		_, err = db.ExecContext(ctx,
			`INSERT INTO plans (name, stripe_price_id, price_cents, interval, features, sort_order) VALUES ($1, $2, $3, $4, $5, $6)`,
			p.Name, p.StripePriceID, p.PriceCents, p.Interval, featuresJSON, p.SortOrder,
		)
		if err != nil {
			log.Fatalf("inserting plan '%s': %v", p.Name, err)
		}
		fmt.Printf("Plan '%s' created ($%d/mo)\n", p.Name, p.PriceCents/100)
	}
}
