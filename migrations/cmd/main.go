package main

import (
	"context"
	"fmt"
	"log"
	"os"
	_ "tempgalias/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	// Database connection string
	godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")

	// Connect to the database using pgx
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Ping the database to check if the connection is working
	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := goose.SetDialect("pgx"); err != nil {
		panic(err)
	}

	db := stdlib.OpenDBFromPool(pool)
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	// Get the Goose migrations directory
	migrationsDir := "./migrations"

	// Run migrations
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run main.go [up|down|status|create]")
	}

	command := os.Args[1]

	if err := goose.RunContext(context.Background(), command, db, migrationsDir, os.Args[2:]...); err != nil {
		log.Fatalf("Failed to run goose command: %v", err)
	}

	fmt.Println("Migrations completed successfully.")
}
