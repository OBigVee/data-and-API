package main

import "time"

type Profile struct {
	ID string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Gender string `db:"gender" json:"gender"`
	GenderProbability  float64   `db:"gender_probability" json:"gender_probability"`
	SampleSize         int       `db:"sample_size" json:"sample_size"`
	Age                int       `db:"age" json:"age"`
	AgeGroup           string    `db:"age_group" json:"age_group"`
	CountryID          string    `db:"country_id" json:"country_id"`
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
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Gender    string `json:"gender" db:"gender"`
	Age       int    `json:"age" db:"age"`
	AgeGroup  string `json:"age_group" db:"age_group"`
	CountryID string `json:"country_id" db:"country_id"`
}