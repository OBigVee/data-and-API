// This file handles reading the seed_profiles.json file
// and inserting the data into the database it generates the required UUID v7 IDs
// performs a batch insert operation

package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func SeedDatabase(db *sqlx.DB) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM profiles")
	if err == nil && count > 0 {
		log.Printf("Database already seeded with %d profiles, skipping seeding.", count)
		return
	}
	// 1. Read the JSON file(seed_profiles.json)
	content, err := os.ReadFile("seed_profiles.json")
	if err != nil {
		log.Printf("CRTICAL: Could not find seed_profiles.json file: %v", err)
		return
	}

	// Unmarshal the JSON data into Profile structs
	// The JSON file has a top-level "profiles" key that contains an array of profiles, so we  todefine a wrapper struct
	var profilesWrapper struct {
		Profiles []Profile `json:"profiles"`
	}
	if err := json.Unmarshal(content, &profilesWrapper); err != nil {
		log.Printf("Error parsing seed JSON: %v", err)
		return
	}

	// generate UUID v7 IDs for each profile 
	for i := range profilesWrapper.Profiles {
		id, _ := uuid.NewV7()
		profilesWrapper.Profiles[i].ID = id.String()
		profilesWrapper.Profiles[i].CreatedAt = time.Now().UTC()
	}

	// batch insert profiles into the database using NamedExec
	query := `INSERT INTO profiles (id, name, gender, gender_probability, age, age_group, country_id, country_name, country_probability, created_at) 
				VALUES (:id, :name, :gender, :gender_probability, :age, :age_group, :country_id, :country_name, :country_probability, :created_at)
				ON CONFLICT (name) DO NOTHING`
	_, err = db.NamedExec(query, profilesWrapper.Profiles)
	if err != nil {
		log.Printf("Seeding failed: %v", err)
		return
	}
	
	log.Println("Successfully seeded 2026 profiles from seed_profiles.json")
}

