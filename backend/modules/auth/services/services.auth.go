package services

import (
	"backend/modules/auth/helpers/apperrors"
	"backend/modules/auth/models"

	"context"
	"log"
	"math/rand"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

// AccountService acts as a struct for injecting an implementation of
// EmployeeRepository for use in the service methods
type employeeService struct {
	EmployeeRepository models.EmployeeRepository
}

// ConfigAccountService will hold repositories that will eventually be injected
// into this service layer
type ConfigAccountService struct {
	EmployeeRepository models.EmployeeRepository
}

// NewAccountService is a factory function for initializing an employeeService
// with it's repository layer dependencies
func NewEmployeeService(ctx *ConfigAccountService) models.EmployeeService {
	return &employeeService{
		EmployeeRepository: ctx.EmployeeRepository,
	}
}

var argon2Params = &argon2id.Params{
	Memory:      62500,
	Iterations:  2,
	Parallelism: 16,
	SaltLength:  32,
	KeyLength:   64,
}

/*********Employee actions*********/

// Signup verifies that the account does not already exist and
// creates the account as an employee
func (service *employeeService) Signup(ctx context.Context, accountData *models.EmployeeAccountData, headersData *models.SigninPayload) error {
	// If an account field is empty throw an error
	if accountData.Email == "" {
		log.Printf("Service Layer Error: email is null")
		return apperrors.NewInternal()
	}

	// If the account email is valid (not already in the DB) then continue

	hashedPassword, err := argon2id.CreateHash(accountData.Password, argon2Params)
	if err != nil {
		log.Panicf("Unable to encrypt accountData password from new Signup with given email: %v\n", accountData.Email)
		return apperrors.NewInternal()
	}
	accountData.Password = hashedPassword //The hashed password is saved to the DB instead of the plain text. User Signin will compare the client's input password to the stored hash and validity

	accountData.IDCode = uuid.New()

	accountData.CreatedAt = time.Now().UTC()
	accountData.UpdatedAt = time.Now().UTC()
	accountData.EmployeeIdentityData.UpdatedAt = time.Now().UTC()

	if err := service.EmployeeRepository.CreateAccount(ctx, accountData, headersData); err != nil {
		return err
	}

	return nil
}

/*
Signin REFACTOR: instead of taking in a models, taking an Email and Password, and return an EmployeeAccountData model.
Pass in email and password params not that we are not mutating the payload.
*/
func (service *employeeService) Signin(ctx context.Context, payload *models.SigninPayload) (*models.EmployeeAccountData, error) {
	account, err := service.EmployeeRepository.FindAccountByEmail(ctx, payload.Email)
	if err != nil {
		return nil, apperrors.NewAuthorization("Debug 1 Invalid email or password")
	}

	// authMatch, err := service.signinAuth(ctx, payload)
	// if err != nil {
	// 	return nil, err
	// }
	// if !authMatch {
	// 	return nil, err
	// }

	// Verify the password against the hash
	accountMatch, err := argon2id.ComparePasswordAndHash(payload.Password, account.Password)
	if err != nil {
		log.Printf("Error: Could not ComparePasswordAndHash successfully")
		return nil, apperrors.NewInternal()
	}
	if !accountMatch {
		return nil, apperrors.NewAuthorization("Debug 2 Invalid email or password boolean")
	}

	account.Password = ""

	//*payload = *account
	return account, nil
}

// GetUser returns a single user based on their uuid
func (service *employeeService) GetAccount(ctx context.Context, accountID uuid.UUID) (*models.EmployeeAccountData, error) {
	log.Print("Hello GetAccount & param accountID:  ", accountID)
	account, err := service.EmployeeRepository.GetAccountDataForClientProfile(ctx, accountID)
	if err != nil {
		log.Printf("error with GetAccount service. Error: %v", err)
		return nil, err
	}
	log.Print("Hit end of GetAccount")
	return account, nil
}

func (service *employeeService) GetJWTValues(ctx context.Context, email string) (*models.EmployeeAccountData, error) {
	fetchJWTValues, err := service.EmployeeRepository.GetJWTValues(ctx, email)
	return fetchJWTValues, err
}

func (service *employeeService) UpdateAccount(ctx context.Context, account *models.EmployeeAccountData) error {
	err := service.EmployeeRepository.UpdateAccount(ctx, account)
	if err != nil {
		return err
	}

	return nil
}

func (service *employeeService) DeleteAccount(ctx context.Context, accountID string, deleteAccountID uuid.UUID) error {
	err := service.EmployeeRepository.DeleteAccount(ctx, deleteAccountID)
	if err != nil {
		log.Printf("erorr when trying to delete account: %v, Error: %v", deleteAccountID, err)
		return err
	}

	return nil
}

// generatePhoneNumber outputs a random 10 digit phone number. Will improve upon later.
func generatePhoneNumber() string {
	var numbers = []rune("1234567890")
	var n int = 10
	b := make([]rune, n)
	for i := range b {
		b[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(b)
}

/*********Security Funcs*********/
func (service *employeeService) signinAuth(ctx context.Context, reqPayload *models.SigninPayload) (bool, error) {
	accountMatch, err := service.EmployeeRepository.SigninHistoryAuth(ctx, reqPayload)
	if err != nil {
		log.Print("error finding account")
		return false, err
	}
	return accountMatch, nil
}
