package models

import (
	"time"
)

// Colony struct models the fields for the Extrasolar Colonies. Colony government level is
// Planet, Country, Province, District, Municipality
type Colony struct {
	Name           string
	Capital        string
	FoundationDate string
	Population     int
	// Countries map[countries]string
	// Provinces map[provinces]string
	// Regions map[regions]string
	// District map[districts]string
	CreatedAt time.Time
	UpdatedAt time.Time
}
