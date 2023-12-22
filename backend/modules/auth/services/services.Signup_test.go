package services

import (
	"backend/internal/helpers/apperrors"
	"backend/internal/mocks"
	"backend/internal/models"
	"time"

	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignup(t *testing.T) {
	t.Run("Success - Simple Signup", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserData := &models.EmployeeAccountData{
			Email:    "crag@tarr.gov",
			Password: "1234567890-=",
		}

		mockEmployeeRepository := new(mocks.MockEmployeeRepository)
		accountService := NewEmployeeService(&ConfigAccountService{
			EmployeeRepository: mockEmployeeRepository,
		})

		// New test setup: use the Run method to modify the passed in data
		// when testing CreateUser. We can then have the Return method
		// return a no error
		mockEmployeeRepository.On("CreateAccount", mock.AnythingOfType("*context.emptyCtx"), mockUserData).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*models.EmployeeAccountData)
				userArg.IDCode = uid
			}).Return(nil)

		ctx := context.TODO()
		err := accountService.Signup(ctx, mockUserData)

		assert.NoError(t, err)

		// assert that the mockUserData has it's IDCode field populated
		assert.Equal(t, uid, mockUserData.IDCode)

		mockEmployeeRepository.AssertExpectations(t)
	})

	t.Run("Success - Full Signup", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockAccountData := &models.EmployeeAccountData{
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

		mockEmployeeRepository := new(mocks.MockEmployeeRepository)
		accountService := NewEmployeeService(&ConfigAccountService{
			EmployeeRepository: mockEmployeeRepository,
		})

		// New test setup: use the Run method to modify the passed in data
		// when testing CreateUser. We can then have the Return method
		// return a no error
		mockEmployeeRepository.On("CreateAccount", mock.AnythingOfType("*context.emptyCtx"), mockAccountData).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*models.EmployeeAccountData)
				userArg.IDCode = uid
			}).Return(nil)

		ctx := context.TODO()
		err := accountService.Signup(ctx, mockAccountData)

		assert.NoError(t, err)

		// assert that the mockUserData has it's IDCode field populated
		assert.Equal(t, uid, mockAccountData.IDCode)

		mockEmployeeRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUserData := &models.EmployeeAccountData{
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

		mockEmployeeRepository := new(mocks.MockEmployeeRepository)

		accountService := NewEmployeeService(&ConfigAccountService{
			EmployeeRepository: mockEmployeeRepository,
		})

		mockError := apperrors.NewConflict("email", mockUserData.Email)

		mockEmployeeRepository.On("CreateAccount", mock.AnythingOfType("*context.emptyCtx"), mockUserData).
			Return(mockError)

		ctx := context.TODO()
		err := accountService.Signup(ctx, mockUserData)

		assert.EqualError(t, err, mockError.Error())
		mockEmployeeRepository.AssertExpectations(t)
	})
}
