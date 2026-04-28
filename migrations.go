package main

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// RunMigrations creates all required database tables if they don't already exist.
// This is idempotent — safe to run on every startup.
func RunMigrations(db *sqlx.DB) {
	log.Println("Running database migrations...")

	// Profiles table (Stage 2 — already exists, but safe to re-run)
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS profiles (
			id VARCHAR PRIMARY KEY,
			name VARCHAR UNIQUE NOT NULL,
			gender VARCHAR,
			gender_probability DOUBLE PRECISION,
			age INTEGER,
			age_group VARCHAR,
			country_id VARCHAR,
			country_name VARCHAR,
			country_probability DOUBLE PRECISION,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create profiles table: %v", err)
	}

	// Users table (Stage 3)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR PRIMARY KEY,
			github_id VARCHAR UNIQUE NOT NULL,
			username VARCHAR NOT NULL,
			email VARCHAR DEFAULT '',
			avatar_url VARCHAR DEFAULT '',
			role VARCHAR DEFAULT 'analyst' CHECK (role IN ('admin', 'analyst')),
			is_active BOOLEAN DEFAULT TRUE,
			last_login_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Refresh tokens table (Stage 3)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id VARCHAR PRIMARY KEY,
			user_id VARCHAR NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token_hash VARCHAR NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create refresh_tokens table: %v", err)
	}

	log.Println("Database migrations completed successfully.")
}
