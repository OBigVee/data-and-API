package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func main() {
	// Load .env file (ignore error if not present — production uses real env vars)
	godotenv.Load()

	var err error
	dsn := os.Getenv("DATABASE_URL")
	log.Println("Connecting to database...")
	db, err = sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations (idempotent — safe on every startup)
	RunMigrations(db)

	// Seed profiles data
	SeedDatabase(db)

	// Start background cleanup goroutines
	CleanupExpiredStates()
	CleanupExpiredTokens()

	r := chi.NewRouter()

	// Global middleware
	r.Use(RequestLogger)
	r.Use(CORSMiddleware)

	// ──────────────────────────────────────────
	// Auth routes (public, rate limited 10/min)
	// ──────────────────────────────────────────
	r.Group(func(r chi.Router) {
		r.Use(RateLimitMiddleware(authRateLimiter, 10, time.Minute))

		r.Get("/auth/github", GitHubLoginHandler)
		r.Get("/auth/github/callback", GitHubCallbackHandler)
		r.Post("/auth/refresh", RefreshTokenHandler)
		r.Post("/auth/logout", LogoutHandler)
	})

	// ──────────────────────────────────────────
	// Protected API routes
	// Auth + API version + CSRF + rate limit 60/min
	// ──────────────────────────────────────────
	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		r.Use(CSRFProtection)
		r.Use(RequireAPIVersion)
		r.Use(RateLimitMiddleware(apiRateLimiter, 60, time.Minute))

		// Read-only endpoints (admin + analyst)
		r.Get("/api/profiles", ListProfilesHandler)
		r.Get("/api/profiles/search", SearchProfilesHandler)
		r.Get("/api/profiles/export", ExportProfilesHandler)
		r.Get("/api/profiles/{id}", GetSingleProfileHandler)

		// Write endpoints (admin only)
		r.With(RequireRole("admin")).Post("/api/profiles", CreateProfileHandler)
		r.With(RequireRole("admin")).Delete("/api/profiles/{id}", DeleteProfileHandler)
	})

	// ──────────────────────────────────────────
	// Authenticated user info (requires auth, no API version)
	// ──────────────────────────────────────────
	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		r.Get("/auth/me", MeHandler)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
