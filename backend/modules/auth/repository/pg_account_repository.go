package repository

import (
	"backend/modules/auth/helpers/apperrors"
	"backend/modules/auth/models"

	"context"
	"fmt"
	"log"
	"time"

	//"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	//"github.com/lib/pq"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	//"github.com/rs/xid"
)

type pgAccountRepository struct {
	DB *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) models.EmployeeRepository {
	return &pgAccountRepository{
		DB: db,
	}
}

/******************** Core CRUD Functions ********************/

// CreateUser adds the user account to the PG DB from the client data
func (pgRepo *pgAccountRepository) CreateAccount(ctx context.Context, account *models.EmployeeAccountData, headersData *models.SigninPayload) error {
	queryAccountData := `INSERT INTO accounts_employee (id_code, email, password, phone_number, job_title, office_address, security_access_level, language, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *`
	queryIdentityData := `INSERT INTO accounts_identity (id_code, first_name, middle_name, last_name, sex, gender, age, height, home_address, birthdate, birthplace, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING *`
	log.Printf("Hello from account_repo layer: %v", account)
	// try writing this pgRepo.DB.GetContext(ctx, account, queryAccountData, account) to reduce manual account.field
	if err := pgRepo.DB.GetContext(ctx, account, queryAccountData, account.IDCode, account.Email, account.Password, account.PhoneNumber, account.JobTitle,
		account.OfficeAddress, account.SecurityAccessLevel, account.Language, account.CreatedAt, account.UpdatedAt); err != nil {
		// check the unique constraint
		if err, ok := err.(*pgconn.PgError); ok && err.Code == "unique_violation" {
			log.Printf("error1: Could not create an account with email: %v, or phone number: %v. Reason: %v\n", account.Email, account.PhoneNumber, err.Code)
			return apperrors.NewConflict("email", account.Email)
		}
		log.Printf("error2: Could not create an account with email: %v. Reason: %v\n", account.Email, err)
		return apperrors.NewInternal()
	}
	log.Printf("Debug: accounts_employee Error checks passed. Review insert for correctness.\n")
	log.Printf("Data passed: id_code: %v, email:%v, password:%v, phone_number:%v, job_title:%v, office_address:%v, secuirty_access_level: %v, langauge: %v\n, created_at:%v, updated_at: %v",
		account.IDCode, account.Email, account.Password, account.PhoneNumber, account.JobTitle, account.OfficeAddress, account.SecurityAccessLevel, account.Language, account.CreatedAt, account.UpdatedAt)

	if err := pgRepo.DB.GetContext(ctx, account, queryIdentityData, account.IDCode, account.EmployeeIdentityData.FirstName, account.EmployeeIdentityData.MiddleName,
		account.EmployeeIdentityData.LastName, account.EmployeeIdentityData.Sex, account.EmployeeIdentityData.Gender, account.EmployeeIdentityData.Age, account.EmployeeIdentityData.Height, account.EmployeeIdentityData.HomeAddress,
		account.EmployeeIdentityData.Birthdate, account.EmployeeIdentityData.Birthplace, account.EmployeeIdentityData.UpdatedAt); err != nil {
		log.Printf("error 3: Could not insert data into accounts_identity. Reason: %v", err)
		return apperrors.NewInternal()
	}
	log.Printf("Debug: accounts_identity error checks passed. Review inserts for correctness")
	log.Printf("Data Passed: id_code:%v, first_name:%v, middle_name:%v, last_name:%v, sex:%v, gender:%v, age:%v, height:%v, home_address:%v, birthdate: %v, birthplace:%v, updated_at: %v\n",
		account.IDCode, account.EmployeeIdentityData.FirstName, account.EmployeeIdentityData.MiddleName, account.EmployeeIdentityData.LastName,
		account.EmployeeIdentityData.Sex, account.EmployeeIdentityData.Gender, account.EmployeeIdentityData.Age, account.EmployeeIdentityData.Height,
		account.EmployeeIdentityData.HomeAddress, account.EmployeeIdentityData.Birthdate, account.EmployeeIdentityData.Birthplace, account.EmployeeIdentityData.UpdatedAt)

	queryHeadersData := `INSERT INTO signin_history (email, ip_address, user_agent, login_timestamp) VALUES ($1, $2, $3, $4)`
	if err := pgRepo.DB.GetContext(ctx, headersData, queryHeadersData, headersData.Email, headersData.Password, headersData.Status, headersData.IPAddress, headersData.UserAgent,
		headersData.LoginTimestamp); err != nil {
		log.Printf("error writing SigninHistory to Postgres. Error:%v\n", err)
	}

	return nil
}

// GetAccountDataForClientProfile returns a complete account profile for the client's account page after logging in successfully.
func (pgRepo *pgAccountRepository) GetAccountDataForClientProfile(ctx context.Context, id uuid.UUID) (*models.EmployeeAccountData, error) {
	employeeData := &models.EmployeeAccountData{}
	identityData := &models.EmployeeIdentityData{}

	queryEmployeeData := `SELECT id_code, email, phone_number, job_title, office_address, security_access_level, language, created_at FROM accounts_employee WHERE id_code = $1`
	queryIdentityData := `SELECT first_name, middle_name, last_name, age, sex, gender, height, home_address, birthplace, birthdate FROM accounts_identity WHERE id_code = $1`

	if err := pgRepo.DB.GetContext(ctx, employeeData, queryEmployeeData, id); err != nil {
		log.Printf("Unable to get account with email: %v. Error:%v\n", id, err)
	}

	if err := pgRepo.DB.GetContext(ctx, identityData, queryIdentityData, id); err != nil {
		log.Printf("Unable to get account identity. Error:%v\n", err)
	}

	account := &models.EmployeeAccountData{
		IDCode:        employeeData.IDCode,
		Email:         employeeData.Email,
		PhoneNumber:   employeeData.PhoneNumber,
		JobTitle:      employeeData.JobTitle,
		OfficeAddress: employeeData.OfficeAddress,
		EmployeeIdentityData: &models.EmployeeIdentityData{
			FirstName:   identityData.FirstName,
			MiddleName:  identityData.MiddleName,
			LastName:    identityData.LastName,
			Age:         identityData.Age,
			Sex:         identityData.Sex,
			Gender:      identityData.Gender,
			Height:      identityData.Height,
			HomeAddress: identityData.HomeAddress,
			Birthplace:  identityData.Birthplace,
			Birthdate:   identityData.Birthdate,
		},
		Language:            employeeData.Language,
		SecurityAccessLevel: employeeData.SecurityAccessLevel,
		CreatedAt:           employeeData.CreatedAt,
	}

	return account, nil
}

func (pgRepos *pgAccountRepository) UpdateAccount(ctx context.Context, account *models.EmployeeAccountData) error {
	query := `UPDATE accounts_employee 
						SET FieldsWeAreUpdating name=:name email=:email 
						WHERE id_code=:id_code
						RETURNING *;
	`
	// nstmt is short for Prepared Named Statment
	prepNameStmt, err := pgRepos.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		log.Printf("Unable to prepare account update using query: %v\n", err)
		return apperrors.NewInternal()
	}

	if err := prepNameStmt.GetContext(ctx, account, account); err != nil {
		log.Printf("Failed to update account details for account: %v\n", account)
		return apperrors.NewInternal()
	}

	return nil
}

/******************** Sys Admin Controls  ********************/

// _DeleteAccount is a SysAdmin command that removes the account from the active accounts tables.
// In the future we might not remove, instead lock it out and maintain the account for records keeping.
func (pgRepo *pgAccountRepository) DeleteAccount(ctx context.Context, account uuid.UUID) error {
	query := `DELETE from accounts_employee WHERE id_code = $1`

	result, err := pgRepo.DB.ExecContext(ctx, query, account)
	if err != nil {
		log.Printf("Error when deleting account with id_code: %v", account)
		return err
	}

	rowAffected, err := result.RowsAffected()
	if rowAffected == 0 {
		log.Printf("Could not delete account with ID Code: %v. Account not found.", account)
		return err
	}

	query = `DELETE FROM accounts_identity WHERE id_code = $1`
	result, err = pgRepo.DB.ExecContext(ctx, query, account)
	if err != nil {
		log.Printf("Error when deleting account with id_code: %v", account)
		return err
	}
	rowAffected, err = result.RowsAffected()
	if rowAffected == 0 {
		log.Printf("Could not delete account with ID Code: %v. Account not found.", account)
		return err
	}

	return nil
}

/******************** Search Functions ********************/

// FindAccountByEmail searches the DB for the account that matches the email parameter and returns the account.
func (pgRepo *pgAccountRepository) FindAccountByEmail(ctx context.Context, email string) (*models.EmployeeAccountData, error) {
	accountData := &models.EmployeeAccountData{
		EmployeeIdentityData: &models.EmployeeIdentityData{},
		//SecurityAccessLevel:  &models.SecurityAccessLevel{},
	}
	query := "SELECT id_code, email, password FROM accounts_employee WHERE email = $1"
	//query := "SELECT * FROM accounts_employee WHERE email = $1"
	if err := pgRepo.DB.GetContext(ctx, accountData, query, email); err != nil {
		log.Printf("Failed to get account with email: %v. Error:%v\n", email, err)
		return accountData, apperrors.NewNotFound("email", email)
	}
	return accountData, nil
}

// FindByID Query's the database to return a user based on the ID passed into it
func (pgRepo *pgAccountRepository) FindAccountDataByID(ctx context.Context, idcode uuid.UUID) (*models.EmployeeAccountData, error) {
	account := &models.EmployeeAccountData{}
	log.Print("Hello pgRepo param idcode ", idcode)

	query := "SELECT * FROM accounts_employee WHERE id_code = $1"

	// Check for errors as the error could be other than "Not Found"
	if err := pgRepo.DB.GetContext(ctx, account, query, idcode); err != nil {
		return account, apperrors.NewNotFound("idcode", idcode.String())
	}
	log.Print(account)

	return account, nil
}

func (pgRepo *pgAccountRepository) FindIdentityDataByID(ctx context.Context, idcode uuid.UUID) (*models.EmployeeIdentityData, error) {
	account := &models.EmployeeIdentityData{}
	log.Print("Hello pgRepo param idcode ", idcode)

	query := "SELECT * FROM accounts_identity WHERE id_code = $1"

	// Check for errors as the error could be other than "Not Found"
	if err := pgRepo.DB.GetContext(ctx, account, query, idcode); err != nil {
		return account, apperrors.NewNotFound("idcode", idcode.String())
	}
	return account, nil
}

func (pgRepo *pgAccountRepository) FindExistingPhoneNumber(ctx context.Context, phoneNumber string) (bool, error) {
	query := `SELECT * FROM accounts_data WHERE phone_number = $1`

	if err := pgRepo.DB.QueryRowContext(ctx, query, phoneNumber); err != nil {
		return false, fmt.Errorf("error: Phone Number %v already exists", phoneNumber)
	}
	return true, nil
}

// Not used. Leave here for possible future use of finely crafted JWT values
func (pgRepo *pgAccountRepository) GetJWTValues(ctx context.Context, email string) (*models.EmployeeAccountData, error) {
	accountData := &models.EmployeeAccountData{}
	accountIdentity := &models.EmployeeIdentityData{}
	queryAccountData := "SELECT id_code, email, password, phone_number, job_title, office_address FROM accounts_employee WHERE email = $1"
	queryIdentityData := "SELECT * FROM accounts_identity WHERE id_code = $1"
	if err := pgRepo.DB.GetContext(ctx, accountData, queryAccountData, email); err != nil {
		log.Printf("Unable to get account with email: %v. Error:%v\n", email, err)
	}

	id_code := accountData.IDCode
	if err := pgRepo.DB.GetContext(ctx, accountIdentity, queryIdentityData, id_code); err != nil {
		log.Printf("Unable to get account identity. Error:%v\n", err)
	}

	jwt := &models.EmployeeAccountData{
		IDCode:        accountData.IDCode,
		Email:         accountData.Email,
		JobTitle:      accountData.JobTitle,
		OfficeAddress: accountData.OfficeAddress,
		EmployeeIdentityData: &models.EmployeeIdentityData{
			FirstName: accountIdentity.FirstName,
			LastName:  accountIdentity.LastName,
		},
		SecurityAccessLevel: accountData.SecurityAccessLevel,
		//SecurityAccessLevel: &models.SecurityAccessLevel{},
	}
	return jwt, nil
}

/********** Security Functions **********/

// SigninHistoryAuth checks the signin request headers against those stored in the DB. Later this will make use of cookies.
func (pgRepo *pgAccountRepository) SigninHistoryAuth(ctx context.Context, request *models.SigninPayload) (bool, error) {
	signinHistory := &models.SigninPayload{}

	query := `SELECT DISTINCT ON (signin_count) * FROM signin_history where email = $1`
	if err := pgRepo.DB.GetContext(ctx, signinHistory, query, request.Email); err != nil {
		log.Printf("Unable to get account with email: %v. Error:%v\n", request.Email, err)
	}

	// check early to see if the account has already been locked
	if signinHistory.Locked {
		query := `INSERT INTO signin_history (email, password, status, ip_address, user_agent, login_timestamp) VALUES ($1, $2, $3, $4, $5, $6)`
		if err := pgRepo.DB.GetContext(ctx, request, query, request.Email, request.Password, request.Status, request.IPAddress, request.UserAgent,
			request.LoginTimestamp); err != nil {
			log.Printf("error writing SigninHistory to Postgres. Error:%v\n", err)
		}
		return false, fmt.Errorf("account has been locked")
	}

	// check if the IP Adresses match
	if request.IPAddress != signinHistory.IPAddress {
		request.SigninCount += 1
		request.LoginTimestamp = time.Now().UTC()

		log.Printf("error: account signin IP Address does not match previous signin: %v | request: %v", signinHistory.IPAddress, request.IPAddress)
		query := `INSERT INTO signin_history (email, password, status, signin_count, ip_address, user_agent, login_timestamp) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`
		if err := pgRepo.DB.GetContext(ctx, request, query, request.Email, request.Password, request.Status, request.SigninCount, request.IPAddress,
			request.UserAgent, request.LoginTimestamp); err != nil {
			log.Printf("error writing SigninHistory to Postgres. Error:%v\n", err)
		}

		if signinHistory.SigninCount == 5 {
			log.Printf("error: Too many failed attemps, account %v has been locked", request.Email)
			return false, fmt.Errorf("too many failed attemps, account has been locked")
		}
		//query = `UPDATE `
		// do something with authenticator app or email one time temporary password
		return false, fmt.Errorf("an error occured during signin process")
	}

	// check if the User-Agent headers match
	if request.UserAgent != signinHistory.UserAgent {
		request.SigninCount += 1
		request.LoginTimestamp = time.Now().UTC()

		log.Printf("error: account signin User-Agent does not match previous signin: %v | request: %v", signinHistory.UserAgent, request.UserAgent)
		query := `INSERT INTO signin_history (email, password, status, ip_address, user_agent, login_timestamp) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`
		if err := pgRepo.DB.GetContext(ctx, request, query, request.Email, request.Password, request.Status, request.SigninCount, request.IPAddress, request.UserAgent,
			request.LoginTimestamp); err != nil {
			log.Printf("error writing SigninHistory to Postgres. Error:%v\n", err)
		}

		if signinHistory.SigninCount == 5 {
			signinHistory.Locked = true
			query := `INSTER INTO signing_history (locked) VALUES ($1)`
			if err := pgRepo.DB.GetContext(ctx, request, query, request.Locked); err != nil {
				log.Printf("error writing SigninHistory to Postgres. Error:%v\n", err)
			}
			log.Printf("error: Too many failed attemps, account %v has been locked", request.Email)
			return false, fmt.Errorf("too many failed attemps, account has been locked")
		}

		query = `UPDATE `
		// do something with authenticator app or email one time temporary password
		return false, fmt.Errorf("an error occured during signin process")
	}

	if signinHistory.SigninCount == 5 {

		log.Printf("error: Too many failed attemps, account %v has been locked", request.Email)
		return false, fmt.Errorf("too many failed attemps, account has been locked")
	}
	// Lock the account
	// do some account security send off maybe like a log or alert?
	// ask for a Signin code through email or authenticator app?

	return true, nil

	// fields := reflect.VisibleFields(reflect.TypeOf(struct{ models.SigninPayload }{}))
	// for _, field := range fields {
	// 	fmt.Printf("Key: %s\tType: %s\n", field.Name, field.Type)
	// }
	// return false, nil
}
