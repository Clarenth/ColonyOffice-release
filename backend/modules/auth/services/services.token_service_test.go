package services

import (
	"backend/modules/auth/mocks"
	"backend/modules/auth/models"

	"context"
	"fmt"
	"os"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTokenPairFromUser(t *testing.T) {
	// Hardcode token expiration timers
	var idTokenExpiration int64 = 15 * 60
	var refreshTokenExpiration int64 = 3 * 24 * 2600

	privateKeyFile, _ := os.ReadFile("../../rsa_private_test.pem")
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	publicKeyFile, _ := os.ReadFile("../../rsa_public_test.pem")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)
	refreshSecretKey := "deepstonecrypt" // Make a temp secret key. Add to .env file later on.

	mockTokenRepository := new(mocks.MockTokenRepository)

	tokenService := NewTokenService(&ConfigTokenService{
		TokenRepository:            mockTokenRepository,
		PrivateKey:                 privateKey,
		PublicKey:                  publicKey,
		RefreshSecretKey:           refreshSecretKey,
		IDTokenExpirationSecs:      idTokenExpiration,
		RefreshTokenExpirationSecs: refreshTokenExpiration,
	})

	uid, _ := uuid.NewRandom()
	mockAccountData := &models.EmployeeAccountData{
		IDCode:        uid,
		Email:         "crag@tarr.gov",
		Password:      "1234567890-=",
		PhoneNumber:   "416-555-1234",
		JobTitle:      "Developer",
		OfficeAddress: "62 West Wallaby Street",
		EmployeeIdentityData: &models.EmployeeIdentityData{
			FirstName:     "Crag",
			MiddleName:    "Marr",
			LastName:      "Tarr",
			Age:           32,
			Height:        "165cm",
			Address:       "26 East Wallaby Street",
			Birthplace:    "Paddington, London, Tarrland, Kingdom of Tarrland, Earth",
			Birthdate:     "2200-01-01",
			UpdatedAtDate: time.Now(),
		},
		SecurityAccessLevel: &models.SecurityAccessLevel{
			ClassificationLevel: "Top Secret",
		},
		CreatedAtDate: time.Now(),
		UpdatedAtDate: time.Now(),
	}

	// mock variables for token functions
	accountIDCodeErrorCase, _ := uuid.NewRandom()
	accountErrorCase := &models.EmployeeAccountData{
		IDCode:        accountIDCodeErrorCase,
		Email:         "crag@tarr.gov",
		Password:      "1234567890-=",
		PhoneNumber:   "416-555-1234",
		JobTitle:      "Developer",
		OfficeAddress: "62 West Wallaby Street",
		EmployeeIdentityData: &models.EmployeeIdentityData{
			FirstName:     "Crag",
			MiddleName:    "Marr",
			LastName:      "Tarr",
			Age:           32,
			Height:        "165cm",
			Address:       "26 East Wallaby Street",
			Birthplace:    "Paddington, London, Tarrland, Kingdom of Tarrland, Earth",
			Birthdate:     "2200-01-01",
			UpdatedAtDate: time.Now(),
		},
		SecurityAccessLevel: &models.SecurityAccessLevel{
			ClassificationLevel: "Top Secret",
		},
		CreatedAtDate: time.Now(),
		UpdatedAtDate: time.Now(),
	}
	previousIDToken := "previous_token"

	setSuccessArguments := mock.Arguments{
		mock.AnythingOfType("*context.Emptyctx"),
		mockAccountData.IDCode.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	setErrorArguments := mock.Arguments{
		mock.AnythingOfType("*context.Emptyctx"),
		accountIDCodeErrorCase.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	deleteWithPreviousIDArgs := mock.Arguments{
		mock.AnythingOfType("*context.Emptyctx"),
		mockAccountData.IDCode.String(),
		previousIDToken,
	}

	// mock calls to arguments and responses
	mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)
	mockTokenRepository.On("SetRefreshToken", setErrorArguments...).Return(fmt.Errorf("There was an error with setting the JWT refresh token"))
	mockTokenRepository.On("DeleteRefreshToken", deleteWithPreviousIDArgs...).Return(nil)

	t.Run("Return a token pair with correct values", func(t *testing.T) {
		ctx := context.Background() // was previously TODO(), changed to background because of return type
		tokenPair, err := tokenService.NewTokenPairFromUser(ctx, mockAccountData, previousIDToken)
		assert.NoError(t, err)

		//Call SeRefreshToken with SetRefreshToken
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken sould not be called since previousIDToken is an empty string ""
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteWithPreviousIDArgs...)

		var verifyStringIsTypeString string
		assert.IsType(t, verifyStringIsTypeString, tokenPair.IDToken)

		// Decode the Base64IRL encoded string. It's simpler to use the
		// JWT library already imported
		idTokenClaims := &idTokenCustomClaims{}

		_, err = jwt.ParseWithClaims(tokenPair.IDToken, idTokenClaims, func(t *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})
		assert.NoError(t, err)

		// Assert the claims on idToken
		expectedTokenClaims := []interface{}{
			mockAccountData.IDCode,
			mockAccountData.Email,
			mockAccountData.PhoneNumber,
			mockAccountData.JobTitle,
			mockAccountData.OfficeAddress,
			mockAccountData.EmployeeIdentityData.FirstName,
			mockAccountData.EmployeeIdentityData.MiddleName,
			mockAccountData.EmployeeIdentityData.LastName,
			mockAccountData.EmployeeIdentityData.Age,
			mockAccountData.EmployeeIdentityData.Height,
			mockAccountData.EmployeeIdentityData.Address,
			mockAccountData.EmployeeIdentityData.Birthdate,
			mockAccountData.EmployeeIdentityData.Birthplace,
			mockAccountData.SecurityAccessLevel.ClassificationLevel,
		}
		actualTokenClaims := []interface{}{
			idTokenClaims.IDCode,
			idTokenClaims.Email,
			idTokenClaims.PhoneNumber,
			idTokenClaims.JobTitle,
			idTokenClaims.OfficeAddress,
			idTokenClaims.FirstName,
			idTokenClaims.LastName,
			idTokenClaims.UpdatedAtDate,
		}

		assert.ElementsMatch(t, expectedTokenClaims, actualTokenClaims)
		//assert.Empty(t, idTokenClaims.Account.Password) // Never pass password through JSON

		expiresAt := time.Unix(idTokenClaims.ExpiresAt.Unix(), 0)
		expectedExpiresAt := time.Now().Add(time.Duration(idTokenExpiration) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		refreshTokenClaims := &refreshTokenPayload{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken, refreshTokenClaims, func(t *jwt.Token) (interface{}, error) {
			return []byte(refreshSecretKey), nil
		})

		assert.IsType(t, verifyStringIsTypeString, tokenPair.RefreshToken)

		// assert the claims on the refresh token
		assert.NoError(t, err)
		assert.Equal(t, mockAccountData.IDCode, refreshTokenClaims.IDCode)

		expiresAt = time.Unix(refreshTokenClaims.RegisteredClaims.ExpiresAt.Unix(), 0)
		expectedExpiresAt = time.Now().Add(time.Duration(refreshTokenExpiration) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)
	})

	t.Run("Error setting refresh token", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewTokenPairFromUser(ctx, accountErrorCase, "")
		assert.Error(t, err)

		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setErrorArguments...)

		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Empty string provided for previousID", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewTokenPairFromUser(ctx, mockAccountData, "")
		assert.NoError(t, err)

		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)

		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})
}
