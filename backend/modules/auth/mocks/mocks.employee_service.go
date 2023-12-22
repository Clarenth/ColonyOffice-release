package mocks

import (
	"backend/internal/models"

	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockEmployeeService is a mock type for models.EmployeeService
type MockEmployeeService struct {
	mock.Mock
}

// GetUser is mock of AccountService GetUser
func (m *MockEmployeeService) GetUser(ctx context.Context, uid uuid.UUID) (*models.EmployeeAccountData, error) {
	// args that will be passed to "Return" in the tests, when function
	// is called with a uid. Hence the name "ret"
	ret := m.Called(ctx, uid)

	// first value passed to "Return"
	var r0 *models.EmployeeAccountData
	if ret.Get(0) != nil {
		// we can just return this if we know we won't be passing function to "Return"
		r0 = ret.Get(0).(*models.EmployeeAccountData)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// Signup is a mock of AccountService Signup
func (m *MockEmployeeService) Signup(ctx context.Context, newAccountData *models.EmployeeAccountData) error {
	ret := m.Called(ctx, newAccountData)

	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// Signin is a mock of the Signin AccountService function
func (m *MockEmployeeService) Signin(ctx context.Context, accountCreds *models.EmployeeAccountData) error {
	ret := m.Called(ctx, accountCreds)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}
	return r0
}
