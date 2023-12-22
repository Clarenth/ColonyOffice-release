package handlers

import (
	"backend/modules/auth/models"
)

/*
Schemas are used to structure incoming requests for the server handlers to understand.
Other names include Data Transfer Objects(DTOs) in Java and C#, and payloads or requests in JavaScript.
*/

// signupRequest is used by the Signup handler for incoming client requests
// when creating new accounts
type signupRequest struct {
	Email                string                       `json:"email" binding:"required,email"`
	Password             string                       `json:"password" binding:"required,gte=16,lte=256"`
	PhoneNumber          string                       `json:"phone_number" binding:"required"`
	JobTitle             string                       `json:"job_title" binding:"required"`
	OfficeAddress        string                       `json:"office_address" binding:"required"`
	EmployeeIdentityData *models.EmployeeIdentityData `json:"employee_identity_data" binding:"required"`
	SecurityAccessLevel  string                       `json:"security_access_level" binding:"required"`
	//SecurityAccessLevel  *models.SecurityAccessLevel  `json:"security_access_level" binding:"required"`
}

type updateRequest struct {
}

// signinRequest is used by Signin handler for incoming client requets
// when logging into an account
type signinRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=16,lte=128"`
}

type deleteRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type tokensRequest struct {
	RefreshTokenString string `json:"refreshToken" binding:"required"`
}
