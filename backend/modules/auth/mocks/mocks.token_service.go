package mocks

import (
	"backend/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) NewTokenPairFromUser(ctx context.Context, account *models.EmployeeAccountData, prevTokenID string) (*models.TokenPair, error) {
	ret := m.Called(ctx, account, prevTokenID)

	// first value passed to "Return"
	var r0 *models.TokenPair
	if ret.Get(0) != nil {
		//return if we know we won't be passing a function to Return
		r0 = ret.Get(0).(*models.TokenPair)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockTokenService) ValidateIDToken(tokenString string) (*models.JWTToken, error) {
	ret := m.Called(tokenString)

	var r0 *models.JWTToken
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.JWTToken)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
