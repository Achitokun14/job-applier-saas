package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Usage: go run cmd/migrate/main.go up|down|version|force N
//
// Environment:
//   DATABASE_URL - PostgreSQL connection string (e.g. postgres://user:pass@host:5432/dbname?sslmode=disable)
//
// Commands:
//   up      - Run all pending migrations
//   down    - Rollback the last migration
//   version - Show current migration version
//   force N - Force set the migration version to N (use to fix dirty state)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "file://migrations"
	}

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	command := os.Args[1]

	switch command {
	case "up":
		log.Println("Running all pending migrations...")
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No pending migrations to apply.")
				return
			}
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migrations applied successfully.")

	case "down":
		log.Println("Rolling back the last migration...")
		if err := m.Steps(-1); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No migrations to rollback.")
				return
			}
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Rollback completed successfully.")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				fmt.Println("No migrations have been applied yet.")
				return
			}
			log.Fatalf("Failed to get version: %v", err)
		}
		fmt.Printf("Current migration version: %d (dirty: %v)\n", version, dirty)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("force command requires a version number: force N")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		log.Printf("Forcing migration version to %d...\n", version)
		if err := m.Force(version); err != nil {
			log.Fatalf("Force version failed: %v", err)
		}
		log.Printf("Migration version forced to %d.\n", version)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run cmd/migrate/main.go <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up        Run all pending migrations")
	fmt.Println("  down      Rollback the last migration")
	fmt.Println("  version   Show current migration version")
	fmt.Println("  force N   Force set the migration version to N")
	fmt.Println()
	fmt.Println("Environment:")
	fmt.Println("  DATABASE_URL      PostgreSQL connection string (required)")
	fmt.Println("  MIGRATIONS_PATH   Path to migrations directory (default: file://migrations)")
}
