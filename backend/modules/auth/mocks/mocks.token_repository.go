package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) SetRefreshToken(ctx context.Context, accountID string, tokenID string, tokenExpiresIn time.Duration) error {
	ret := m.Called(ctx, accountID, tokenExpiresIn)

	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m *MockTokenRepository) DeleteRefreshToken(ctx context.Context, accountID string, tokenExpiresIn string) error {
	ret := m.Called(ctx, accountID, tokenExpiresIn)

	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
