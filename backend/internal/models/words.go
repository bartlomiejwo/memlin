package models

type Word struct {
	ID            int     `json:"id,omitempty"`
	Word          string  `json:"word"`
	Language      string  `json:"language"`
	Pronunciation string  `json:"pronunciation"`
	Category      string  `json:"category"`
	Level         string  `json:"level"`
	Popularity    float64 `json:"popularity"`
}

type WordCreateParams struct {
	LanguageCode  string  `json:"language_code"`
	Word          string  `json:"word"`
	Pronunciation string  `json:"pronunciation"`
	Category      string  `json:"category"`
	Level         string  `json:"level"`
	Popularity    float64 `json:"popularity"`
}
