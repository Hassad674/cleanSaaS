package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/hassad/boilerplateSaaS/backend/internal/config"
)

func main() {
	cfg := config.Load()

	flag.Parse()
	command := flag.Arg(0)

	if command == "" {
		fmt.Println("Usage: go run cmd/migrate/main.go <command>")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  up        Apply all pending migrations")
		fmt.Println("  down      Rollback the last migration")
		fmt.Println("  down-all  Rollback ALL migrations")
		fmt.Println("  status    Show current migration version")
		fmt.Println("  force N   Force set version N (use after a failed migration)")
		os.Exit(1)
	}

	m, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up failed: %v", err)
		}
		version, dirty, _ := m.Version()
		fmt.Printf("Migrations applied. Current version: %d (dirty: %v)\n", version, dirty)

	case "down":
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration down failed: %v", err)
		}
		version, dirty, _ := m.Version()
		fmt.Printf("Rolled back 1 migration. Current version: %d (dirty: %v)\n", version, dirty)

	case "down-all":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration down-all failed: %v", err)
		}
		fmt.Println("All migrations rolled back.")

	case "status":
		version, dirty, err := m.Version()
		if err == migrate.ErrNilVersion {
			fmt.Println("No migrations applied yet.")
		} else if err != nil {
			log.Fatalf("failed to get version: %v", err)
		} else {
			fmt.Printf("Current version: %d (dirty: %v)\n", version, dirty)
		}

	case "force":
		vStr := flag.Arg(1)
		if vStr == "" {
			log.Fatal("force requires a version number: go run cmd/migrate/main.go force 1")
		}
		v, err := strconv.Atoi(vStr)
		if err != nil {
			log.Fatalf("invalid version number: %s", vStr)
		}
		if err := m.Force(v); err != nil {
			log.Fatalf("force failed: %v", err)
		}
		fmt.Printf("Forced version to %d.\n", v)

	default:
		log.Fatalf("unknown command: %s", command)
	}
}
