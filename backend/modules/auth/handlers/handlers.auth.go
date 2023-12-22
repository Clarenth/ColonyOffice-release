package handlers

import (
	"backend/modules/auth/helpers/apperrors"
	"backend/modules/auth/models"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler holds the methods expected by the service layer for handling routing
type handler struct {
	EmployeeService models.EmployeeService
	TokenService    models.TokenService
}

// HandlerConfig holds the server configuration need by by the handler
type HandlerConfig struct {
	EmployeeService models.EmployeeService
	TokenService    models.TokenService
}

func NewAuthHandlers(ctx *HandlerConfig) models.AuthHandlers {
	return &handler{
		EmployeeService: ctx.EmployeeService,
		TokenService:    ctx.TokenService,
	}
}

// CurrentUser handler calls services to get the current users account details
func (handler *handler) CurrentAccount(ctx *gin.Context) {
	// in the future the service will retrieve the user's details through middleware and pass them here.
	accountKey, exists := ctx.Get("account")

	// Errors should not happen as that will be handled by the middleware. This is extra safety handling.
	if !exists {
		log.Printf("Unable to extract request context for unknown user: %v\n", ctx)
		err := apperrors.NewInternal()
		ctx.JSON(err.Status(), gin.H{
			"error on CurrentUser": err,
		})
		return
	}

	employeeIDCode := accountKey.(*models.JWTToken).IDCode
	log.Print(employeeIDCode)

	ctxRequest := ctx.Request.Context()
	getAccountFromServices, err := handler.EmployeeService.GetAccount(ctxRequest, employeeIDCode)
	log.Print(getAccountFromServices.IDCode)

	if err != nil {
		log.Printf("Unable to find account: %v\n%v", employeeIDCode, err)
		e := apperrors.NewNotFound("account", employeeIDCode.String())

		ctx.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"account": getAccountFromServices,
	})
}

func (handler *handler) Signin(ctx *gin.Context) {
	// request takes in a signinRequest sruct to model the incoming JSON request data
	var request signinRequest

	// Bind incoming JSON to a struct and check for validation errors
	if ok := bindData(ctx, &request); !ok {
		return
	}

	payload := &models.SigninPayload{
		Email:          request.Email,
		Password:       request.Password,
		Status:         false,
		IPAddress:      ctx.ClientIP(),
		UserAgent:      ctx.Request.UserAgent(),
		LoginTimestamp: time.Now().UTC(),
	}

	ctxRequest := ctx.Request.Context()
	account, err := handler.EmployeeService.Signin(ctxRequest, payload)
	if err != nil {
		log.Printf("Failed to signin user: %v.\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := handler.TokenService.NewTokenPairFromUser(ctx, account, "")
	if err != nil {
		log.Printf("Error: failed to create tokens: %v\n", err)
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	// Send a status 200 with the authorized tokens
	ctx.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}

func (handler *handler) Signout(ctx *gin.Context) {
	accountID := ctx.MustGet("account")

	ctxRequest := ctx.Request.Context()
	if err := handler.TokenService.Signout(ctxRequest, accountID.(*models.JWTToken).IDCode); err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "account signout successfull",
	})

}

// Signup handler passes user data to create an account for them
func (handler *handler) Signup(ctx *gin.Context) {
	// request will bind the incoming request data from the client
	var request signupRequest

	// Bind incoming JSON to a struct and check for validation errors
	if ok := bindData(ctx, &request); !ok {
		return
	}

	newAccountData := &models.EmployeeAccountData{
		Email:         request.Email,
		Password:      request.Password,
		PhoneNumber:   request.PhoneNumber,
		JobTitle:      request.JobTitle,
		OfficeAddress: request.OfficeAddress,
		EmployeeIdentityData: &models.EmployeeIdentityData{
			FirstName:   request.EmployeeIdentityData.FirstName,
			MiddleName:  request.EmployeeIdentityData.MiddleName,
			LastName:    request.EmployeeIdentityData.LastName,
			Sex:         request.EmployeeIdentityData.Sex,
			Gender:      request.EmployeeIdentityData.Gender,
			Age:         request.EmployeeIdentityData.Age,
			Height:      request.EmployeeIdentityData.Height,
			HomeAddress: request.EmployeeIdentityData.HomeAddress,
			Birthdate:   request.EmployeeIdentityData.Birthdate,
			Birthplace:  request.EmployeeIdentityData.Birthplace,
		},
		SecurityAccessLevel: request.SecurityAccessLevel,
		// SecurityAccessLevel: &models.SecurityAccessLevel{
		// 	ClassificationLevel: request.SecurityAccessLevel.ClassificationLevel,
		// },
		// IPAddress: ctx.ClientIP(),
		// UserAgent: ctx.Request.UserAgent(),
		Language: ctx.Request.Header.Get("Accept-Language"),
	}

	headersData := &models.SigninPayload{
		Email:     request.Email,
		Status:    false,
		IPAddress: ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
	}

	log.Print(newAccountData)
	log.Print(newAccountData.EmployeeIdentityData)
	log.Print(headersData)

	ctxRequest := ctx.Request.Context()
	err := handler.EmployeeService.Signup(ctxRequest, newAccountData, headersData)
	if err != nil {
		log.Printf("Failed to signup account: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := handler.TokenService.NewTokenPairFromUser(ctxRequest, newAccountData, "")
	if err != nil {
		log.Printf("Failed to create tokens for account: %v\n", err.Error())
		// Implement rollback logic that deletes the user DB information if tokens fail
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status_msg": "Account creation successful.",
		"tokens":     tokens,
	})
}

func (handler *handler) UpdateAccount(ctx *gin.Context) {
	authAccount := ctx.MustGet("account").(*models.EmployeeAccountData)

	var request updateRequest

	if ok := bindData(ctx, &request); !ok {
		return
	}

	accountData := &models.EmployeeAccountData{
		IDCode: authAccount.IDCode,
	}

	ctxRequest := ctx.Request.Context()
	err := handler.EmployeeService.UpdateAccount(ctxRequest, accountData)
	if err != nil {
		log.Printf("updates to account %v failed.", err.Error())

		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"account": accountData,
	})
}

func (handler *handler) DeleteAccount(ctx *gin.Context) {
	authHeader := ctx.MustGet("account").(*models.JWTToken)

	id := ctx.Params.ByName("id")
	deleteAccountID := uuid.MustParse(id)

	err := handler.EmployeeService.DeleteAccount(ctx, authHeader.IDCode.String(), deleteAccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"erorr": "Could not delete account",
		})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})

}
