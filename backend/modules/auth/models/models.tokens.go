package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Consider for the future when saving to httponly
// CookieDomain   string
// CookiePath     string
// CookieName     string

// JWTToken struct contains the JWT claims used in system authentication
type JWTToken struct {
	// Account ID Code of the account who owns the JWT
	IDCode uuid.UUID `json:"id_code"`
	// Email belonging to the Account who owns the JWT
	Email string `json:"email"`
	// JobTitle      string    `json:"job_title"`
	// OfficeAddress string    `json:"office_address"`
	// FirstName     string    `json:"first_name"`
	// LastName      string    `json:"last_name"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	// ID           uuid.UUID `json:"_"`
	// UID          uuid.UUID `json:"_"`
	JWTIDCode     uuid.UUID `json:"_"`
	AccountIDCode uuid.UUID `json:"_"`
	SignedString  string    `json:"refreshToken"`
}

// TokenPair is used in a JWT token exchange when a user logs in. IDToken holds the
// unique ID of the JWT. RefreshToken is sent to when the user login's.
type TokenPair struct {
	IDToken
	RefreshToken
}

type IDToken struct {
	SignedString string `json:"idToken"`
}
