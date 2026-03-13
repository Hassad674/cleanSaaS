package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/lib/pq"

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
	seedBlogPosts(ctx, db)

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

func seedBlogPosts(ctx context.Context, db *sql.DB) {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM blog_posts`).Scan(&count); err != nil {
		fmt.Println("Blog posts table not found, skipping blog seed.")
		return
	}
	if count > 0 {
		fmt.Printf("Blog already has %d posts, skipping.\n", count)
		return
	}

	// Get admin user ID for author
	var authorID string
	err := db.QueryRowContext(ctx, `SELECT id FROM users WHERE email = 'admin@cleansaas.dev'`).Scan(&authorID)
	if err != nil {
		fmt.Println("Admin user not found, skipping blog seed.")
		return
	}

	posts := []struct {
		Title           string
		Slug            string
		Excerpt         string
		Content         string
		MetaTitle       string
		MetaDescription string
		Tags            []string
		Status          string
	}{
		{
			Title:   "Getting Started with CleanSaaS",
			Slug:    "getting-started-with-cleansaas",
			Excerpt: "Learn how to set up and customize CleanSaaS for your next project.",
			Content: "# Getting Started with CleanSaaS\n\nCleanSaaS is an open-source boilerplate for building medium-to-large SaaS applications. It provides a solid foundation with authentication, billing, AI chat, file storage, notifications, and more.\n\n## Quick Setup\n\n1. Clone the repository\n2. Run `docker compose up -d` to start PostgreSQL\n3. Apply migrations with `make migrate-up`\n4. Start the backend with `make run`\n5. Start the frontend with `npm run dev`\n\n## Architecture\n\nThe backend follows hexagonal architecture (ports and adapters), ensuring clean separation of concerns. The frontend uses a feature-based structure with Next.js 15.\n\n## Modularity\n\nEvery feature is fully independent and removable. You can use billing without AI, or notifications without storage. The upcoming CLI tool will let you pick only the modules you need.",
			MetaTitle:       "Getting Started with CleanSaaS - Setup Guide",
			MetaDescription: "Step-by-step guide to set up CleanSaaS, the open-source SaaS boilerplate with Go, Next.js, and PostgreSQL.",
			Tags:            []string{"tutorial", "getting-started", "architecture"},
			Status:          "published",
		},
		{
			Title:   "Why We Chose Hexagonal Architecture",
			Slug:    "why-hexagonal-architecture",
			Excerpt: "The reasoning behind our backend architecture choices and how they benefit your project.",
			Content: "# Why We Chose Hexagonal Architecture\n\nWhen building a boilerplate that will be used by thousands of developers, architecture decisions matter more than ever.\n\n## The Problem\n\nMost SaaS starters tightly couple their business logic with their framework, database, and external services. This makes it nearly impossible to swap providers or test in isolation.\n\n## Our Solution\n\nHexagonal architecture (also called ports and adapters) solves this by defining clear boundaries:\n\n- **Domain**: Pure business logic with zero dependencies\n- **Ports**: Interface contracts that define what the system needs\n- **Adapters**: Concrete implementations that can be swapped freely\n- **Application**: Use cases that orchestrate domain and ports\n\n## Real Benefits\n\n- **Swap Stripe for Lemon Squeezy?** Change one adapter file and one line in main.go.\n- **Switch from PostgreSQL to MySQL?** Implement the repository interfaces with a new adapter.\n- **Test business logic without a database?** Mock the repository interface.\n\nThis is not over-engineering — it is the minimum architecture for a system designed to be customized.",
			MetaTitle:       "Hexagonal Architecture in Go - CleanSaaS",
			MetaDescription: "Learn why CleanSaaS uses hexagonal architecture and how it enables easy provider swapping and testing.",
			Tags:            []string{"architecture", "go", "backend"},
			Status:          "published",
		},
		{
			Title:   "Building a Feature-Based Frontend with Next.js",
			Slug:    "feature-based-frontend-nextjs",
			Excerpt: "How our frontend architecture keeps features independent and your codebase maintainable.",
			Content: "# Building a Feature-Based Frontend with Next.js\n\nA feature-based architecture organizes code by business domain rather than technical layer. Here is how CleanSaaS implements this with Next.js 15.\n\n## Structure\n\nEach feature lives in its own folder under src/features/ and contains everything it needs: components, hooks, actions, types, and utilities. Features never import from each other.\n\n## Benefits\n\n- **Remove a feature** by deleting its folder — zero compilation errors elsewhere\n- **Add a feature** by creating a new folder — no modifications to existing code\n- **Understand a feature** by reading one directory — everything is co-located\n\n## Composition\n\nPages in the app/ directory are thin composition layers that import from features and combine them. This is the only place where features meet.\n\n## Design Tokens\n\nWe use CSS custom properties for theming. No hardcoded Tailwind colors — everything goes through semantic tokens like bg-card, text-primary, and border-border.",
			MetaTitle:       "Feature-Based Frontend Architecture with Next.js 15",
			MetaDescription: "Discover the feature-based frontend architecture in CleanSaaS with Next.js 15, design tokens, and full modularity.",
			Tags:            []string{"frontend", "nextjs", "architecture"},
			Status:          "published",
		},
	}

	for _, p := range posts {
		_, err := db.ExecContext(ctx,
			`INSERT INTO blog_posts (author_id, title, slug, excerpt, content, meta_title, meta_description, tags, status, published_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())`,
			authorID, p.Title, p.Slug, p.Excerpt, p.Content, p.MetaTitle, p.MetaDescription, pq.Array(p.Tags), p.Status,
		)
		if err != nil {
			log.Fatalf("inserting blog post '%s': %v", p.Title, err)
		}
		fmt.Printf("Blog post created: '%s'\n", p.Title)
	}
}
