package services

import (
	"backend/internal/mocks"
	"backend/internal/models"

	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockEmployeeResponse := &models.EmployeeAccountData{
			IDCode: uid,
			Email:  "crag@tarr.gov",
			EmployeeIdentityData: &models.EmployeeIdentityData{
				FirstName: "Crag Tarr",
			},
		}

		mockEmployeeRepository := new(mocks.MockEmployeeRepository)
		employeeService := NewEmployeeService(&ConfigAccountService{
			EmployeeRepository: mockEmployeeRepository,
		})
		mockEmployeeRepository.On("FindByID", mock.Anything, uid).Return(mockEmployeeResponse, nil)

		ctx := context.TODO()
		user, err := employeeService.GetUser(ctx, uid)

		assert.NoError(t, err)
		assert.Equal(t, user, mockEmployeeResponse)
		mockEmployeeRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockEmoloyeeRepository := new(mocks.MockEmployeeRepository)
		employeeRepository := NewEmployeeService(&ConfigAccountService{
			EmployeeRepository: mockEmoloyeeRepository,
		})

		mockEmoloyeeRepository.On("FindByID", mock.Anything, uid).Return(nil, fmt.Errorf("An error occured down the call chain"))

		ctx := context.TODO()
		user, err := employeeRepository.GetUser(ctx, uid)

		assert.Nil(t, user)
		assert.Error(t, err)
		mockEmoloyeeRepository.AssertExpectations(t)
	})
}
