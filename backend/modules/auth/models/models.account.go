package models

import (
	"time"

	"github.com/google/uuid"
)

// EmployeeAccountData models the Account's assignable communications fields
// and the Account's secuirty privillege and clearance
type EmployeeAccountData struct {
	// IDCode struct tag is a 32 digit Universally unique identifier that
	// exist, is unique/exclusive/reserved, etc. for each employee
	IDCode uuid.UUID `db:"id_code" json:"id_code"` //binding:"required
	// Email is a unique email for each employee of the Colony Office
	Email string `db:"email" json:"email"`
	// Password is a string between 16 and 256 characters long
	Password string `db:"password" json:"-"`
	// Phone Number is the contact phone number issued to each employee of the Colony Office
	PhoneNumber string `db:"phone_number" json:"phone_number"`
	// Job Title is the current role the employee is assigned
	JobTitle string `db:"job_title" json:"job_title"`
	// OfficeAddress is the physical location of where the employee is assigned to work at.
	OfficeAddress string `db:"office_address" json:"office_address"`
	// EmployeeIdentityData is a pointer to a limited-mutabale identity data.
	// This should be change only with employee consent, and official documentation
	// e.g. Name change, sex transition, etc.
	EmployeeIdentityData *EmployeeIdentityData `db:"employee_identity_data" json:"employee_identity_data"`
	// SecurityAccessLevel is the account security level allowing access to higher levels of system access
	SecurityAccessLevel string `db:"security_access_level" json:"security_access_level"`
	//SecurityAccessLevel *SecurityAccessLevel `db:"classification_level" json:"classification_level"`
	// Language is the language set in the 'Accept-Language' header during signup
	Language string `db:"language" json:"-"`
	//CreationDate is the timestamp that the user account was created at
	CreatedAt time.Time `db:"created_at" json:"employment_date"`
	// UpdatedAtDate is the timestamp for when the EmployeeAccountData was updated at any point in time
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

// EmployeeIdentityData models the Employee's personal information
type EmployeeIdentityData struct {
	IDCode uuid.UUID `db:"id_code" json:"-"`
	// FirstName struct field is for the First Name of the user
	FirstName string `db:"first_name" json:"first_name"`
	// MiddleName struct field is for the Middle Name of the user. *Optionally leave
	// blank if user has no Middle Name
	MiddleName string `db:"middle_name" json:"middle_name"`
	// LastName struct field is for the Last Name of the user
	LastName string `db:"last_name" json:"last_name"`
	// Sex struct field is the biological sex of the employee
	Sex string `db:"sex" json:"sex"`
	// Gender struct field is the gender identity of the employee
	Gender string `db:"gender" json:"gender"`
	// Age struct field is the current age of the employee
	Age int `db:"age" json:",string"`
	// Height is represented in metric using centimetres
	Height string `db:"height" json:"height"`
	// Address is the physical location of where the employee is living at.
	// Mail and first priority contact should be sent here
	HomeAddress string `db:"home_address" json:"home_address"`
	// Birthplace *Colony `json:"birth_place"`
	// Birthdate struct field is the date of when the user was born. Format is in
	// ISO 8601 Year-Month-Day YYYY-MM-DD
	Birthdate string `db:"birthdate" json:"birthdate"`
	// Birthplace struct field is the physical location of where the user was born.
	// Format is in Municipality, District, Province, Administrative Territory, Planet
	Birthplace string `db:"birthplace" json:"birthplace"`
	// UpdatedAtDate is the timestamp for when the EmployeeIdentityData was updated at any point in time
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

// Separate SecurityAccessLevel and User Created Documents fields into their own
// tables for higher secuirty protection when performing look ups
// From Surton: access level and document data is encrypted into a separate table
// split data once more for data that’s privileged and encrypt that portion
// ? Maybe. Need to plan this out
// Do not forget encryption, seperation of privillege based on Security Level
