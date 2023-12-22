package services

import (
	"backend/modules/auth/models"
	"fmt"

	"crypto/rsa"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

/*
jwt_tokens.go handles the functions called by services.token_service.go when
working with JWTs such as generating a token ID,  or refreshing a token.
Modifying the values of a JWT should be done here instead of inside token service.

Investigate the use of JWE - JSON Web Encryption, or JWT encryption
*/

// idTokenCustomClaims struct contains the JWT claims used in Signin authentication
type idTokenCustomClaims struct {
	IDCode        uuid.UUID `json:"id_code"`
	Email         string    `json:"email"`
	JobTitle      string    `json:"job_title"`
	OfficeAddress string    `json:"office_address"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	jwt.RegisteredClaims
}

// generateIDToken returns an IDToken
func generateIDToken(account *models.EmployeeAccountData, key *rsa.PrivateKey, tokenExpiresAt int64) (string, error) {
	tokenIssuedAt := time.Now()                                                          // Begin Lifecycle at current time of creation
	tokenLifecycleTime := tokenIssuedAt.Add(time.Duration(tokenExpiresAt) * time.Second) // Lifecycle is 15 minutes from time of issue
	log.Print(tokenIssuedAt)
	log.Print(tokenLifecycleTime)

	claims := models.JWTToken{
		IDCode: account.IDCode,
		Email:  account.Email,
		// JobTitle:      account.JobTitle,
		// OfficeAddress: account.OfficeAddress,
		// FirstName:     account.EmployeeIdentityData.FirstName,
		// LastName:      account.EmployeeIdentityData.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(tokenIssuedAt), // tokenIssuedAt
			ExpiresAt: jwt.NewNumericDate(tokenLifecycleTime),
		},
		/*
			later, plan for the use of private claims. These can include UpdatedAtDate: account.UpdatedAtDate, maybe the hashed & salted password
		*/
	}
	log.Print(claims.RegisteredClaims.IssuedAt)
	log.Print(claims.RegisteredClaims.ExpiresAt)

	newJWTToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenSignedString, err := newJWTToken.SignedString(key)

	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return tokenSignedString, nil
}

// refreshTokenData holds the signed JWT and it's ID. ID is returned to avoid re-parsing the signed string
type refreshTokenData struct {
	SignedString   string
	ID             uuid.UUID
	ExpirationTime time.Duration
}

// RefreshTokenPayload structures the payload sent from the frontend user when
// requesting a refresh of their token. We can extract the account IDCode for use in
// other operations of the application (Redis)
type refreshTokenPayload struct {
	IDCode uuid.UUID `json:"id_code"`
	jwt.RegisteredClaims
}

// generateRefreshToken creates a refresh token that stores the account ID
func generateRefreshToken(uid uuid.UUID, key string, tokenExpiresAt int64) (*refreshTokenData, error) {
	currentTime := time.Now()
	tokenExpiration := currentTime.Add(time.Duration(tokenExpiresAt) * time.Second)
	tokenID, err := uuid.NewRandom()

	if err != nil {
		log.Println("Error generating Refresh Token ID")
		return nil, err
	}

	claims := refreshTokenPayload{
		IDCode: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: currentTime},
			ExpiresAt: &jwt.NumericDate{Time: tokenExpiration},
			ID:        tokenID.String(),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := newToken.SignedString([]byte(key))
	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}

	return &refreshTokenData{
		SignedString:   signedString,
		ID:             tokenID,
		ExpirationTime: tokenExpiration.Sub(currentTime),
	}, nil
}

func validateIDToken(tokenString string, key *rsa.PublicKey) (*models.JWTToken, error) {
	claims := &models.JWTToken{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// If there is an error then we'll return and handle it's logging in the service layer
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token ID is invalid")
	}

	claims, ok := token.Claims.(*models.JWTToken)
	if !ok {
		return nil, fmt.Errorf("token ID is valid, but could not parse claims")
	}

	return claims, nil
}

// validateRefreshToken uses the JWT secret key to validate a refresh token
func validateRefreshToken(tokenString string, key string) (*refreshTokenPayload, error) {
	claims := &refreshTokenPayload{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	// For now return the error to the service layer and log the result there.
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("debug jwt_tokens.go: Refresh token is invalid")
	}

	claims, ok := token.Claims.(*refreshTokenPayload)
	if !ok {
		return nil, fmt.Errorf("refresh token was valid, but could not parse the claims")
	}

	return claims, nil
}
