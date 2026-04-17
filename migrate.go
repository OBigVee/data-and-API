package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Using the connection string you provided
	connStr := "postgresql://neondb_owner:npg_zwv20OhTuKWH@ep-sparkling-credit-ana2ljln-pooler.c-6.us-east-1.aws.neon.tech/neondb?sslmode=require"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	query := `CREATE TABLE IF NOT EXISTS profiles (
		id UUID PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		gender VARCHAR(50),
		gender_probability FLOAT,
		sample_size INT,
		age INT,
		age_group VARCHAR(50),
		country_id VARCHAR(10),
		country_probability FLOAT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	log.Println("Schema initialized successfully!")
}
