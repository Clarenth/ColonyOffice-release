package services

import (
	"backend/modules/auth/helpers/apperrors"
	"backend/modules/auth/models"

	"context"
	"crypto/rsa"
	"log"

	"github.com/google/uuid"
)

/*
In the future implement the encryption of tokens being passed to hide the JWT claims.
*/

// TokenService models the token request data passed to the service layer
type tokenService struct {
	TokenRepository            models.TokenRepository
	PrivateKey                 *rsa.PrivateKey
	PublicKey                  *rsa.PublicKey
	RefreshSecretKey           string
	IDTokenExpirationSecs      int64
	RefreshTokenExpirationSecs int64
}

// ConfigTokenService will hold the repositories that are injected into
// the service layer
type ConfigTokenService struct {
	TokenRepository            models.TokenRepository
	PrivateKey                 *rsa.PrivateKey
	PublicKey                  *rsa.PublicKey
	RefreshSecretKey           string
	IDTokenExpirationSecs      int64
	RefreshTokenExpirationSecs int64
}

func NewTokenService(config *ConfigTokenService) models.TokenService {
	return &tokenService{
		TokenRepository:            config.TokenRepository,
		PrivateKey:                 config.PrivateKey,
		PublicKey:                  config.PublicKey,
		RefreshSecretKey:           config.RefreshSecretKey,
		IDTokenExpirationSecs:      config.IDTokenExpirationSecs,
		RefreshTokenExpirationSecs: config.RefreshTokenExpirationSecs,
	}
}

// NewTokenPairFromUser creates new tokens for Account signup and signin.
// If a previous token is included it will be removed from the tokens repository.
func (service *tokenService) NewTokenPairFromUser(ctx context.Context, account *models.EmployeeAccountData, previousTokenID string) (*models.TokenPair, error) {
	log.Print("Hello previousTokenID: ", previousTokenID)
	log.Print("Hello id_code ", account.IDCode)
	if previousTokenID != "" {
		if err := service.TokenRepository.DeleteRefreshToken(ctx, account.IDCode.String(), previousTokenID); err != nil {
			log.Printf("Could not delete previous refresh token for account id_code %v, token ID: %v\n", account.IDCode, previousTokenID)

			return nil, err
		}
	}

	// No need for a repository as idToken is unrelated to any data source
	idToken, err := generateIDToken(account, service.PrivateKey, service.IDTokenExpirationSecs)
	if err != nil {
		log.Printf("Error generating new idToken for account ID Code: %v. Error: %v\n", account.IDCode, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(account.IDCode, service.RefreshSecretKey, service.RefreshTokenExpirationSecs)
	if err != nil {
		log.Printf("Error generating refreshToken for account: %v. Error: %v\n", account.IDCode, err.Error())
		return nil, apperrors.NewInternal()
	}

	// On account signup or signin, generate a new token for the session
	if err := service.TokenRepository.SetRefreshToken(ctx, account.IDCode.String(), refreshToken.ID.String(), refreshToken.ExpirationTime); err != nil {
		log.Printf("Error storing tokenID for account id_code: %v. Error: %v\n", account.IDCode, err.Error())
		return nil, apperrors.NewInternal()
	}

	return &models.TokenPair{
		IDToken: models.IDToken{SignedString: idToken},
		RefreshToken: models.RefreshToken{
			// ID:  refreshToken.ID,
			// UID: account.IDCode,
			JWTIDCode:     refreshToken.ID,
			AccountIDCode: account.IDCode,
			SignedString:  refreshToken.SignedString,
		},
	}, nil
}

func (service *tokenService) ValidateIDToken(tokenString string) (*models.JWTToken, error) {
	claims, err := validateIDToken(tokenString, service.PublicKey) // Signs the token using the public RSA Key
	if err != nil {
		log.Printf("Unable to validate or parse ID on Token - Error %v\n: ", err)
		return nil, apperrors.NewAuthorization("Unable to verify ID on token")
	}
	return claims, err
	// return claims.IDCode, err
}

func (service *tokenService) ValidateRefreshToken(tokenString string) (*models.RefreshToken, error) {
	claims, err := validateRefreshToken(tokenString, service.RefreshSecretKey)
	if err != nil {
		log.Printf("unable to validate or parse Refresh Token for token string: %s\n%v\n", tokenString, err)
		return nil, apperrors.NewAuthorization("Unable to verify account from refresh token")
	}

	// JWT standard claims are a string. The model uses UUID for ID so we need to parse it.
	tokenUUID, err := uuid.Parse(claims.ID)
	if err != nil {
		log.Printf("unable to parse claims.IDCode as UUID: %s\n%v\n", claims.IDCode, err)
		return nil, apperrors.NewAuthorization("Unable to verify account from refresh token")
	}

	return &models.RefreshToken{
		// ID:  tokenUUID,
		// UID: claims.IDCode,
		JWTIDCode:     tokenUUID,
		AccountIDCode: claims.IDCode,
		SignedString:  tokenString,
	}, nil
}

func (service *tokenService) Signout(ctx context.Context, accountID uuid.UUID) error {
	return service.TokenRepository.DeleteAccountRefreshTokens(ctx, accountID.String())
}
