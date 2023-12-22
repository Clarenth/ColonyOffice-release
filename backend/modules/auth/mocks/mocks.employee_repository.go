package mocks

import (
	"backend/modules/auth/models"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockEmployeeRepository is a mock type of the
// models interface models.EmployeeRepository
type MockEmployeeRepository struct {
	mock.Mock
}

// FindByID is a mock of the EmployeeRepository method FindByID
func (m *MockEmployeeRepository) FindByID(ctx context.Context, uid uuid.UUID) (*models.EmployeeAccountData, error) {
	ret := m.Called(ctx, uid)

	var r0 *models.EmployeeAccountData

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.EmployeeAccountData)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockEmployeeRepository) CreateAccount(ctx context.Context, user *models.EmployeeAccountData) error {
	ret := m.Called(ctx, user)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m *MockEmployeeRepository) FindAccountByEmail(ctx context.Context, accountEmail string) (*models.EmployeeAccountData, error) {
	ret := m.Called(ctx, accountEmail)

	var r0 *models.EmployeeAccountData
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.EmployeeAccountData)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
