package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	//"github.com/shoenig/test/interfaces"
)

var countryCodeToName = map[string]string{
    // Africa
    "NG": "Nigeria", "GH": "Ghana", "KE": "Kenya", "ZA": "South Africa",
    "TZ": "Tanzania", "UG": "Uganda", "RW": "Rwanda", "CM": "Cameroon",
    "CI": "Ivory Coast", "SN": "Senegal", "ET": "Ethiopia", "EG": "Egypt",
    "DZ": "Algeria", "MA": "Morocco", "TN": "Tunisia", "SD": "Sudan",
    "LY": "Libya", "MR": "Mauritania", "ML": "Mali", "NE": "Niger",
    "BJ": "Benin", "TG": "Togo", "BF": "Burkina Faso", "GW": "Guinea-Bissau",
    "GN": "Guinea", "SL": "Sierra Leone", "LR": "Liberia", "ZW": "Zimbabwe",
    "MW": "Malawi", "MZ": "Mozambique", "ZM": "Zambia", "NA": "Namibia",
    "BW": "Botswana", "LS": "Lesotho", "SZ": "Eswatini", "MU": "Mauritius",
    "SC": "Seychelles", "KM": "Comoros",

    // Americas
    "CA": "Canada", "US": "United States", "MX": "Mexico", "BR": "Brazil",
    "AR": "Argentina", "CO": "Colombia", "PE": "Peru", "VE": "Venezuela",
    "CL": "Chile", "EC": "Ecuador", "BO": "Bolivia", "PY": "Paraguay",
    "UY": "Uruguay", "CR": "Costa Rica", "PA": "Panama", "DO": "Dominican Republic",
    "HT": "Haiti", "CU": "Cuba", "JM": "Jamaica", "BS": "Bahamas",
    "BB": "Barbados", "GD": "Grenada", "AG": "Antigua and Barbuda",
    "DM": "Dominica", "LC": "Saint Lucia", "VC": "Saint Vincent and the Grenadines",
    "KN": "Saint Kitts and Nevis", "TT": "Trinidad and Tobago", "SR": "Suriname",
    "GY": "Guyana", "PR": "Puerto Rico", "VI": "U.S. Virgin Islands",

    // Europe
    "HU": "Hungary", "FR": "France", "DE": "Germany", "GB": "United Kingdom",
    "IT": "Italy", "ES": "Spain", "NL": "Netherlands", "BE": "Belgium",
    "GR": "Greece", "PT": "Portugal", "CZ": "Czech Republic", "SK": "Slovakia",
    "PL": "Poland", "UA": "Ukraine", "RO": "Romania", "IE": "Ireland",
    "AT": "Austria", "CH": "Switzerland", "NO": "Norway", "SE": "Sweden",
    "FI": "Finland", "DK": "Denmark", "IS": "Iceland", "TR": "Turkey",
    "RU": "Russia",

    // Asia
    "JP": "Japan", "CN": "China", "IN": "India", "TH": "Thailand",
    "MY": "Malaysia", "SG": "Singapore", "ID": "Indonesia", "PH": "Philippines",
    "VN": "Vietnam", "KH": "Cambodia",

    // Oceania & Territories
    "AU": "Australia", "NZ": "New Zealand", "FJ": "Fiji", "PG": "Papua New Guinea",
    "SB": "Solomon Islands", "VU": "Vanuatu", "WS": "Samoa", "TO": "Tonga",
    "TV": "Tuvalu", "KI": "Kiribati", "NR": "Nauru", "PW": "Palau",
    "FM": "Micronesia", "MH": "Marshall Islands", "CK": "Cook Islands",
    "NU": "Niue", "PF": "French Polynesia", "NC": "New Caledonia",
    "WF": "Wallis and Futuna", "TK": "Tokelau", "PN": "Pitcairn Islands",
    "AS": "American Samoa", "GU": "Guam", "MP": "Northern Mariana Islands",
}


func Send502(w http.ResponseWriter, api string) {
	w.WriteHeader(http.StatusBadGateway)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "error",
		"message": fmt.Sprintf("%s returned an invalid response", api),
	})
}

func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status: "error",
		Message: message,
	})
}

// NL search Handler
func SearchProfilesHandler(w http.ResponseWriter, r *http.Request) {
	rawQuery := strings.ToLower(r.URL.Query().Get("q"))

	if rawQuery == "" {
		sendError(w, http.StatusBadRequest, "Missing or empty parameter")
		return
	}

	interpreted := r.URL.Query()

	// Logic: Gender inference based on keywords in the query
	if strings.Contains(rawQuery, "male") || strings.Contains(rawQuery, "men"){
		interpreted.Set("gender", "male")
	} else if strings.Contains(rawQuery, "female") || strings.Contains(rawQuery, "women"){
		interpreted.Set("gender", "female")
	}

	// Logic: Age group inference based on keywords in the query
	if strings.Contains(rawQuery, "child") || strings.Contains(rawQuery, "kid") || strings.Contains(rawQuery, "children"){
		interpreted.Set("age_group", "child")
	} else if strings.Contains(rawQuery, "teen") || strings.Contains(rawQuery, "adolescent"){
		interpreted.Set("age_group", "teenager")
	} else if strings.Contains(rawQuery, "adult") {
		interpreted.Set("age_group", "adult")
	} else if strings.Contains(rawQuery, "senior") || strings.Contains(rawQuery, "elderly") {
		interpreted.Set("age_group", "senior")
	}

	// Logic: "young" mapping (16-24)
	if strings.Contains(rawQuery, "young") {
		interpreted.Set("min_age", "16")
		interpreted.Set("max_age", "24")
	}

	// Logic : "Above X" using Regex to extract age value and set min_age filter
	reAbove := regexp.MustCompile(`above\s+(\d+)`)
	if matches := reAbove.FindStringSubmatch(rawQuery); len(matches) > 1 {
		interpreted.Set("min_age", matches[1])
	}

	// Logic: Country inference based on keywords in the query
	for code, name := range countryCodeToName {
		// regex looks for whole word matches to avoid false positives (eg. "US" in "music")
		pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(strings.ToLower(name)))
		matched, _ := regexp.MatchString(pattern, rawQuery)

		if matched {
			interpreted.Set("country_id", code)
			break
		}
	}

	// If we could not infer any filters from the query, return an error
	// if len(interpreted) <= 1 {
	// 	sendError(w, http.StatusOK, "Could not interpret query")
	// 	return
	// }

	// pass the interpreted query to the main list logic
	r.URL.RawQuery = interpreted.Encode()
	ListProfilesHandler(w, r)
}


// main list handler (filters + sorting + pagination) - this is the main handler that the NL search handler will call after interpreting the query and setting the appropriate filters

func ListProfilesHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	
	query := "SELECT id, name, gender, gender_probability, age, age_group, country_id, country_name, country_probability, created_at FROM profiles WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM profiles WHERE 1=1"
	var args []interface{}
	argID := 1

	// Dynamic filter Builder
	addFilter := func(field string, operator string, value string) {
		if value != "" {
			condition := fmt.Sprintf(" AND %s %s $%d", field, operator, argID)
			query += condition
			countQuery += condition
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
	if sortBy != "age"	 && sortBy != "gender_probability" { sortBy = "created_at" }
	order := strings.ToUpper(q.Get("order"))
	if order != "ASC" { order = "DESC" }
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)
	
	// Pagination
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 || limit > 50 { limit = 10 }
	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 { page = 1 }
	offset := (page - 1) * limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// Execute
	var profiles []Profile
	if err := db.Select(&profiles, query, args...); err != nil { 
		//log.Printf("SQL Error: %v", err)
		sendError(w, http.StatusInternalServerError, "Database error")
		return
	}
	
	var total int
	db.Get(&total, countQuery, args...)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProfileResponse{
		Status: "success",
		Page: page,
		Limit: limit,
		Total: total,
		Data: profiles,
	})
}


func CreateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	nameIface, ok := input["name"]
	if !ok || nameIface == nil {
		sendError(w, http.StatusBadRequest, "Missing or empty name")
		return
	}

	nameStr, ok := nameIface.(string)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Invalid type",
		})
		return
	}

	name := strings.ToLower(strings.TrimSpace(nameStr))
	if name == "" {
		sendError(w, http.StatusBadRequest, "Name cannot be empty")
		return
	}

	// Idempotency check - if profile already exists, return it instead of creating a new one
	var existing Profile
	err := db.Get(&existing, "SELECT * FROM profiles WHERE name = $1", name)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"data":    existing,
		})
		return
	}

	// concurrent API CALLS
	var gData GenderizeResp
	var aData AgifyResp
	var nData NationalizeResp
	var gErr, aErr, nErr error
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		gErr = fetchJSON(fmt.Sprintf("https://api.genderize.io?name=%s", name), &gData)
	}()

	go func() {
		defer wg.Done()
		aErr = fetchJSON(fmt.Sprintf("https://api.agify.io?name=%s", name), &aData)
	}()

	go func() {
		defer wg.Done()
		nErr = fetchJSON(fmt.Sprintf("https://api.nationalize.io?name=%s", name), &nData)
	}()

	wg.Wait()

	if gErr != nil || gData.Gender == "" || gData.Count == 0 {
		Send502(w, "Genderize")
		return
	}
	if aErr != nil || aData.Age == nil {
		Send502(w, "Agify")
		return
	}
	if nErr != nil || len(nData.Country) == 0 {
		Send502(w, "Nationalize")
		return
	}

	// classification logic
	topCountry := nData.Country[0]
	for _, c := range nData.Country {
		if c.Probability > topCountry.Probability {
			topCountry = c
		}
	}

	ageGroup := "senior"
	ageValue := *aData.Age
	if ageValue <= 12 {
		ageGroup = "child"
	} else if ageValue <= 19 {
		ageGroup = "teenager"
	} else if ageValue <= 59 {
		ageGroup = "adult"
	}

	newID, _ := uuid.NewV7()

	countryName := countryCodeToName[topCountry.CountryID]
	if countryName == "" {
		countryName = "Unknown"
	}

	profile := Profile{
		ID:                 newID.String(),
		Name:               name,
		Gender:             gData.Gender,
		GenderProbability:  gData.Probability,
		Age:                ageValue,
		AgeGroup:           ageGroup,
		CountryID:          topCountry.CountryID,
		CountryName:   countryName,
		CountryProbability: topCountry.Probability,
		CreatedAt:          time.Now().UTC(),
	}

	

	// Save to DB
	query := `INSERT INTO profiles (id, name, gender, gender_probability, age, age_group, country_id, country_name, country_probability, created_at)
	VALUES (:id, :name, :gender, :gender_probability, :age, :age_group, :country_id, :country_name, :country_probability, :created_at)`

	_, err = db.NamedExec(query, &profile)

	if err != nil {
		log.Printf("Database Insert Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Failed to save profile",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   profile,
	})
}



// Get singe profile by ID
func GetSingleProfileHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var profile Profile
	err := db.Get(&profile, "SELECT * FROM profiles WHERE id = $1", id)
	if err != nil {
		sendError(w, http.StatusNotFound, "Profile not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   profile,
	})
}

// List all profiles with filtering

func ListProfilesHandlerDeprecated(w http.ResponseWriter, r *http.Request) {
	gender := r.URL.Query().Get("gender")
	country := r.URL.Query().Get("country_id")
	ageGroup := r.URL.Query().Get("age_group")

	query := "SELECT id, name, gender, age, age_group, country_id FROM profiles WHERE 1=1"
	var args []interface{}
	argCount := 1

	if gender != "" {
		query += fmt.Sprintf(" AND gender ILIKE $%d", argCount)
		args = append(args, gender)
		argCount++
	}
	if country != "" {
		query += fmt.Sprintf(" AND country_id ILIKE $%d", argCount)
		args = append(args, country)
		argCount++
	}
	if ageGroup != "" {
		query += fmt.Sprintf(" AND age_group ILIKE $%d", argCount)
		args = append(args, ageGroup)
		argCount++
	}
	var profiles []ProfileListResp
	err := db.Select(&profiles, query, args...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Failed to fetch profiles",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"count":  len(profiles),
		"data":   profiles,
	})
}

// Delete profile by ID
func DeleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		sendError(w, http.StatusBadRequest, "Missing profile ID")
		return
	}

	result, err := db.Exec("DELETE FROM profiles WHERE id = $1", id)
	if err != nil {
		log.Printf("Delete Error: %v", err)
		//w.WriteHeader(http.StatusInternalServerError)
		sendError(w, http.StatusInternalServerError, "Failed to delete profile")
		
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Error verifying deletion")
		return
	}

	if rowsAffected == 0 {
		sendError(w, http.StatusNotFound, "Profile not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func fetchJSON(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
