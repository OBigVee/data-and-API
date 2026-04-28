package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ExportProfilesHandler exports profiles as CSV with the same filtering/sorting as ListProfilesHandler
func ExportProfilesHandler(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format != "csv" {
		sendError(w, http.StatusBadRequest, "Unsupported export format. Use format=csv")
		return
	}

	q := r.URL.Query()

	query := "SELECT id, name, gender, gender_probability, age, age_group, country_id, country_name, country_probability, created_at FROM profiles WHERE 1=1"
	var args []interface{}
	argID := 1

	// Dynamic filter builder (same as ListProfilesHandler)
	addFilter := func(field string, operator string, value string) {
		if value != "" {
			query += fmt.Sprintf(" AND %s %s $%d", field, operator, argID)
			args = append(args, value)
			argID++
		}
	}

	addFilter("gender", "=", q.Get("gender"))
	addFilter("age_group", "=", q.Get("age_group"))
	addFilter("country_id", "=", q.Get("country_id"))
	addFilter("age", ">=", q.Get("min_age"))
	addFilter("age", "<=", q.Get("max_age"))
	addFilter("gender_probability", ">=", q.Get("min_gender_probability"))
	addFilter("country_probability", ">=", q.Get("min_country_probability"))

	// Sorting
	sortBy := q.Get("sort_by")
	if sortBy != "" {
		if sortBy != "age" && sortBy != "gender_probability" && sortBy != "created_at" {
			sendError(w, http.StatusBadRequest, "Invalid sort field")
			return
		}
	} else {
		sortBy = "created_at"
	}

	order := strings.ToUpper(q.Get("order"))
	if order != "ASC" {
		order = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)

	// Execute (no pagination — export all matching)
	var profiles []Profile
	if err := db.Select(&profiles, query, args...); err != nil {
		sendError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Generate CSV
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("profiles_%s.csv", timestamp)

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Header row
	writer.Write([]string{
		"id", "name", "gender", "gender_probability", "age", "age_group",
		"country_id", "country_name", "country_probability", "created_at",
	})

	// Data rows
	for _, p := range profiles {
		writer.Write([]string{
			p.ID,
			p.Name,
			p.Gender,
			fmt.Sprintf("%.2f", p.GenderProbability),
			strconv.Itoa(p.Age),
			p.AgeGroup,
			p.CountryID,
			p.CountryName,
			fmt.Sprintf("%.2f", p.CountryProbability),
			p.CreatedAt.Format(time.RFC3339),
		})
	}
}
