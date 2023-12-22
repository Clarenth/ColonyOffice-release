package models

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandlers interface {
	CurrentAccount(ctx *gin.Context)
	Signin(ctx *gin.Context)
	Signout(ctx *gin.Context)
	Signup(ctx *gin.Context)
	UpdateAccount(ctx *gin.Context)
	Tokens(ctx *gin.Context)
	DeleteAccount(ctx *gin.Context)
}

type TokenHandlers interface {
	Tokens(ctx *gin.Context)
}

// EmployeeService implements methods the methods handler layer expects
// any service it interacts with to implement
type EmployeeService interface {
	GetAccount(ctx context.Context, uid uuid.UUID) (*EmployeeAccountData, error)
	GetJWTValues(ctx context.Context, email string) (*EmployeeAccountData, error)
	Signin(ctx context.Context, payload *SigninPayload) (*EmployeeAccountData, error)
	Signup(ctx context.Context, newAccountData *EmployeeAccountData, headersData *SigninPayload) error
	UpdateAccount(ctx context.Context, account *EmployeeAccountData) error
	DeleteAccount(ctx context.Context, accountID string, deleteAccountID uuid.UUID) error
}

// EmployeeRepository defines the methods the service layer expects
// any repository it interacts with to implement
type EmployeeRepository interface {
	CreateAccount(ctx context.Context, accountData *EmployeeAccountData, headersData *SigninPayload) error
	FindAccountByEmail(ctx context.Context, email string) (*EmployeeAccountData, error)
	FindAccountDataByID(ctx context.Context, uid uuid.UUID) (*EmployeeAccountData, error)
	FindExistingPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)
	FindIdentityDataByID(ctx context.Context, uid uuid.UUID) (*EmployeeIdentityData, error)
	GetAccountDataForClientProfile(ctx context.Context, idcode uuid.UUID) (*EmployeeAccountData, error)
	GetJWTValues(ctx context.Context, email string) (*EmployeeAccountData, error)
	SigninHistoryAuth(ctx context.Context, reqPayload *SigninPayload) (bool, error)
	UpdateAccount(ctx context.Context, account *EmployeeAccountData) error
	DeleteAccount(ctx context.Context, deleteAccountID uuid.UUID) error
}

// TokenRepository defines the methods used when accounts perform actions like signin, logout, data requests, etc.
type TokenRepository interface {
	SetRefreshToken(ctx context.Context, accountID string, tokenID string, tokenExpireTime time.Duration) error
	DeleteAccountRefreshTokens(ctx context.Context, accountID string) error
	DeleteRefreshToken(ctx context.Context, accountID string, tokenID string) error
}

// TokenService defines the methods used to generate account tokens on account actions
type TokenService interface {
	NewTokenPairFromUser(ctx context.Context, account *EmployeeAccountData, previousTokenID string) (*TokenPair, error)
	ValidateIDToken(tokenString string) (*JWTToken, error)
	ValidateRefreshToken(tokenString string) (*RefreshToken, error)
	Signout(ctx context.Context, accountID uuid.UUID) error
}
