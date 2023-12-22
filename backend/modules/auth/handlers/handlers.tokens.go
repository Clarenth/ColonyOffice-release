package handlers

import (
	"backend/modules/auth/helpers/apperrors"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *handler) Tokens(ctx *gin.Context) {
	var request tokensRequest

	if ok := bindData(ctx, &request); !ok {
		return
	}

	log.Print("Hello from Tokens request", request)

	ctxRequest := ctx.Request.Context()

	// Verify the refresh of the JWT
	refreshToken, err := handler.TokenService.ValidateRefreshToken(request.RefreshTokenString)
	if err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Get the account the refresh token belongs to.
	//account, err := handler.EmployeeService.GetAccount(ctxRequest, refreshToken.UID)
	account, err := handler.EmployeeService.GetAccount(ctxRequest, refreshToken.AccountIDCode)
	if err != nil {
		log.Print("Hit error with GetAccount")
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Create a new pair of tokens
	//tokens, err := handler.TokenService.NewTokenPairFromUser(ctxRequest, account, refreshToken.ID.String())
	tokens, err := handler.TokenService.NewTokenPairFromUser(ctxRequest, account, refreshToken.JWTIDCode.String())
	if err != nil {
		log.Printf("Failed to create tokens for account %+v. Error: %v\n", account, err.Error())

		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
	}

	//handler.TokenService.InvalidateOldRefreshToken(ctx, request.RefreshTokenString)
	log.Print("Hello new tokens", tokens)
	ctx.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})

}
