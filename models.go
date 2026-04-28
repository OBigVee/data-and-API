package main

import "time"

// ──────────────────────────────────────────────
// Profile models (Stage 2 — unchanged)
// ──────────────────────────────────────────────

type Profile struct {
	ID                 string    `db:"id" json:"id"`
	Name               string    `db:"name" json:"name"`
	Gender             string    `db:"gender" json:"gender"`
	GenderProbability  float64   `db:"gender_probability" json:"gender_probability"`
	Age                int       `db:"age" json:"age"`
	AgeGroup           string    `db:"age_group" json:"age_group"`
	CountryID          string    `db:"country_id" json:"country_id"`
	CountryName        string    `db:"country_name" json:"country_name"`
	CountryProbability float64   `db:"country_probability" json:"country_probability"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
}

type GenderizeResp struct {
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
}

type NationalizeResp struct {
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

type AgifyResp struct {
	Age *int `json:"age"`
}

type ProfileListResp struct {
	ID                 string    `json:"id" db:"id"`
	Name               string    `json:"name" db:"name"`
	Gender             string    `json:"gender" db:"gender"`
	Age                int       `json:"age" db:"age"`
	AgeGroup           string    `json:"age_group" db:"age_group"`
	CountryID          string    `json:"country_id" db:"country_id"`
	CountryName        string    `json:"country_name" db:"country_name"`
	CountryProbability float64   `json:"country_probability" db:"country_probability"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

// ──────────────────────────────────────────────
// Updated pagination response (Stage 3)
// ──────────────────────────────────────────────

type PaginationLinks struct {
	Self string  `json:"self"`
	Next *string `json:"next"`
	Prev *string `json:"prev"`
}

type ProfileResponse struct {
	Status     string          `json:"status"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	Total      int             `json:"total"`
	TotalPages int             `json:"total_pages"`
	Links      PaginationLinks `json:"links"`
	Data       []Profile       `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ──────────────────────────────────────────────
// User models (Stage 3)
// ──────────────────────────────────────────────

type User struct {
	ID          string     `db:"id" json:"id"`
	GitHubID    string     `db:"github_id" json:"github_id"`
	Username    string     `db:"username" json:"username"`
	Email       string     `db:"email" json:"email"`
	AvatarURL   string     `db:"avatar_url" json:"avatar_url"`
	Role        string     `db:"role" json:"role"`
	IsActive    bool       `db:"is_active" json:"is_active"`
	LastLoginAt *time.Time `db:"last_login_at" json:"last_login_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

type RefreshToken struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}