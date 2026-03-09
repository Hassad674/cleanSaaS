package main

import (
	"log"
	"net/http"

	"github.com/hassad/boilerplateSaaS/backend/internal/adapter/postgres"
	appauth "github.com/hassad/boilerplateSaaS/backend/internal/app/auth"
	appuser "github.com/hassad/boilerplateSaaS/backend/internal/app/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/config"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
)

func main() {
	cfg := config.Load()

	// Database
	db := postgres.NewDB(cfg.DatabaseURL)
	defer db.Close()

	// Repositories
	userRepo := postgres.NewUserRepository(db)

	// JWT
	jwtMaker := jwt.NewMaker(cfg.JWTSecret)

	// App services
	authSvc := appauth.NewService(userRepo, nil, jwtMaker) // email service nil for now
	userSvc := appuser.NewService(userRepo)

	// Router
	router := handler.NewRouter(authSvc, userSvc, cfg.JWTSecret)

	// Start
	log.Printf("API server starting on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
