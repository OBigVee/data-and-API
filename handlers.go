package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	//"github.com/shoenig/test/interfaces"
)

func send502(w http.ResponseWriter, api string) {
	w.WriteHeader(http.StatusBadGateway)
	json.NewEncoder(w).Encode(map[string]string{ 
		"status" : "error",
		"message": fmt.Sprintf("%s returned an invalid response", api),
	})
}

func CreateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Missing or empty name",
		})
		return
	}

	nameIface, ok := input["name"]
	if !ok || nameIface == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Missing or empty name",
		})
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Missing or empty name",
		})
		return
	}
	
	// Idempotency check - if profile already exists, return it instead of creating a new one
	var existing Profile
	err := db.Get(&existing, "SELECT * FROM profiles WHERE name = $1", name)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"message": "Profile already exists", 
			"data": existing,
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
		send502(w, "Genderize")
		return
	}
	if aErr != nil || aData.Age == nil {
		send502(w, "Agify")
		return
	}
	if nErr != nil || len(nData.Country) == 0 {
		send502(w, "Nationalize")
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

	newID, _:= uuid.NewV7()

	profile := Profile{
		ID: newID.String(),
		Name: name,
		Gender: gData.Gender,
		GenderProbability: gData.Probability,
		SampleSize: gData.Count,
		Age: ageValue,
		AgeGroup: ageGroup,
		CountryID: topCountry.CountryID,
		CountryProbability: topCountry.Probability,
		CreatedAt: time.Now().UTC(),
	}


	// Save to DB
	_, err = db.NamedExec(`INSERT INTO profiles (id, name, gender, gender_probability, sample_size, age, age_group, country_id, country_probability, created_at)
	VALUES (:id, :name, :gender, :gender_probability, :sample_size, :age, :age_group, :country_id, :country_probability, :created_at)`, &profile)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "error",
			"message": "Failed to save profile",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": profile,
	})
}

// Get singe profile by ID
func GetSingleProfileHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var profile Profile
	err := db.Get(&profile, "SELECT * FROM profiles WHERE id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "error",
			"message": "Profile not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": profile,
	})
}

// List all profiles with filtering

func ListProfilesHandler(w http.ResponseWriter, r *http.Request) {
	// For simplicity, this example does not implement filtering logic
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
			"status": "error",
			"message": "Failed to fetch profiles",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"count": len(profiles),
		"data": profiles,
	})

}

// Delete profile by ID
func DeleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	result, err := db.Exec("DELETE FROM profiles WHERE id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "error",
			"message": "Failed to delete profile",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "error",
			"message": "Profile not found",
		})
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