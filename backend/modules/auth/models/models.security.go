package models

import "time"

// SigninPayload captures the incoming request's email, password, IP Address, the browser User-Agent, and a timestamp.
// The captured data is initially save during Signup
// and will be saved to the database for comparison during signin requests.
type SigninPayload struct {
	// Email is the email used during Signin
	Email string `db:"email"`
	// Password is the password used during Signin
	Password string `db:"password"`
	// Status is a boolean showing if the signin attempt was successful or not
	Status bool `db:"status"`
	// SigninCount tracks the number of failed signin attemps the user failed at
	SigninCount int `db:"signin_count"`
	// IPAddress is the inital IP Address captured during signup
	IPAddress string `db:"ip_address"`
	// UserAgent is the initial User Agent captured during signup
	UserAgent string `db:"user_agent"`
	// LoginTimestamp is the time that the Signin request was made
	LoginTimestamp time.Time `db:"login_timestamp"`
	Locked         bool      `db:"locked"`
}
