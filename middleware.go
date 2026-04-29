package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ──────────────────────────────────────────────
// Context keys
// ──────────────────────────────────────────────

type contextKey string

const contextKeyUser contextKey = "user"

// ──────────────────────────────────────────────
// RequireAuth — Validates JWT from header or cookie
// ──────────────────────────────────────────────

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Try Authorization header first
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Fall back to HTTP-only cookie
		if tokenString == "" {
			cookie, err := r.Cookie("access_token")
			if err == nil {
				tokenString = cookie.Value
			}
		}

		if tokenString == "" {
			w.Header().Set("Content-Type", "application/json")
			sendError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// Parse and validate JWT
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "insighta-dev-secret-change-me"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			sendError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			sendError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Fetch user from DB to check is_active
		userID, _ := claims["sub"].(string)
		var user User
		err = db.Get(&user, "SELECT * FROM users WHERE id = $1", userID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			sendError(w, http.StatusUnauthorized, "User not found")
			return
		}

		if !user.IsActive {
			w.Header().Set("Content-Type", "application/json")
			sendError(w, http.StatusForbidden, "User account is deactivated")
			return
		}

		// Set user in context
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ──────────────────────────────────────────────
// RequireRole — Checks user role against allowed roles
// ──────────────────────────────────────────────

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(contextKeyUser).(User)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				sendError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			for _, role := range roles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			sendError(w, http.StatusForbidden, "Insufficient permissions")
			return
		})
	}
}

// ──────────────────────────────────────────────
// RequireAPIVersion — Enforces X-API-Version: 1
// ──────────────────────────────────────────────

func RequireAPIVersion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := r.Header.Get("X-API-Version")
		if version == "" {
			w.Header().Set("Content-Type", "application/json")
			sendError(w, http.StatusBadRequest, "API version header required")
			return
		}
		if version != "1" {
			w.Header().Set("Content-Type", "application/json")
			sendError(w, http.StatusBadRequest, "Unsupported API version")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ──────────────────────────────────────────────
// CSRFProtection — Validates CSRF token for web clients
// ──────────────────────────────────────────────

func CSRFProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only enforce on state-changing methods
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// Only enforce for cookie-based auth (web clients)
		if _, err := r.Cookie("access_token"); err != nil {
			// Not a web client — skip CSRF check
			next.ServeHTTP(w, r)
			return
		}

		csrfHeader := r.Header.Get("X-CSRF-Token")
		csrfCookie, err := r.Cookie("csrf_token")

		if err != nil || csrfHeader == "" || csrfHeader != csrfCookie.Value {
			w.Header().Set("Content-Type", "application/json")
			sendError(w, http.StatusForbidden, "Invalid CSRF token")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ──────────────────────────────────────────────
// RateLimiter — In-memory sliding window rate limiter
// ──────────────────────────────────────────────

type rateLimiterEntry struct {
	Timestamps []time.Time
}

type RateLimiterStore struct {
	mu      sync.Mutex
	entries map[string]*rateLimiterEntry
}

func NewRateLimiterStore() *RateLimiterStore {
	return &RateLimiterStore{
		entries: make(map[string]*rateLimiterEntry),
	}
}

var (
	authRateLimiter = NewRateLimiterStore()
	apiRateLimiter  = NewRateLimiterStore()
)

func (s *RateLimiterStore) Allow(key string, limit int, window time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	entry, exists := s.entries[key]
	if !exists {
		entry = &rateLimiterEntry{}
		s.entries[key] = entry
	}

	// Remove timestamps outside the window
	cutoff := now.Add(-window)
	valid := make([]time.Time, 0, len(entry.Timestamps))
	for _, ts := range entry.Timestamps {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}
	entry.Timestamps = valid

	if len(entry.Timestamps) >= limit {
		return false
	}

	entry.Timestamps = append(entry.Timestamps, now)
	return true
}

func RateLimitMiddleware(store *RateLimiterStore, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use user ID if authenticated, otherwise client IP
			key := r.Header.Get("X-Forwarded-For")
			if key == "" {
				key = r.Header.Get("X-Real-IP")
			}
			if key == "" {
				key = r.Header.Get("X-Client-IP")
			}
			
			if key == "" {
				key = r.RemoteAddr
			} else {
				// Handle multiple IPs (take the first one)
				key = strings.TrimSpace(strings.Split(key, ",")[0])
			}

			// Aggressively strip port from whatever key we ended up with
			host, _, err := net.SplitHostPort(key)
			if err == nil {
				key = host
			}

			if user, ok := r.Context().Value(contextKeyUser).(User); ok {
				key = user.ID
			}

			if !store.Allow(key, limit, window) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "60")
				sendError(w, http.StatusTooManyRequests, "Rate limit exceeded. Try again later.")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ──────────────────────────────────────────────
// RequestLogger — Logs method, endpoint, status, response time
// ──────────────────────────────────────────────

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		log.Printf("[%s] %s %d %s", r.Method, r.URL.Path, rec.statusCode, duration)
	})
}

// ──────────────────────────────────────────────
// CORS Middleware — Updated for credentials support
// ──────────────────────────────────────────────

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		webPortalURL := strings.TrimRight(os.Getenv("WEB_PORTAL_URL"), "/")

		// Allow the web portal origin and localhost for development
		allowedOrigins := []string{webPortalURL, "http://localhost:5173", "http://localhost:3000"}
		allowed := false
		for _, o := range allowedOrigins {
			if o != "" && origin == o {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if origin == "" {
			// No origin = same-origin or CLI request
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Version, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-Id")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// ──────────────────────────────────────────────
// CleanupExpiredStates — Background goroutine to clean up old OAuth states
// ──────────────────────────────────────────────

func CleanupExpiredStates() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			oauthStatesMu.Lock()
			for key, entry := range oauthStates {
				if time.Since(entry.CreatedAt) > 10*time.Minute {
					delete(oauthStates, key)
				}
			}
			oauthStatesMu.Unlock()
		}
	}()
}

// ──────────────────────────────────────────────
// CleanupExpiredTokens — Background goroutine to clean up expired refresh tokens
// ──────────────────────────────────────────────

func CleanupExpiredTokens() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			result, err := db.Exec(`DELETE FROM refresh_tokens WHERE expires_at < $1`, time.Now())
			if err != nil {
				log.Printf("Error cleaning expired tokens: %v", err)
			} else {
				rows, _ := result.RowsAffected()
				if rows > 0 {
					log.Printf("Cleaned up %d expired refresh tokens", rows)
				}
			}
		}
	}()
}

// ──────────────────────────────────────────────
// HashToken — Helper to hash tokens for comparison
// ──────────────────────────────────────────────

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
