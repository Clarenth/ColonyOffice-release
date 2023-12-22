package middleware

import (
	"backend/modules/auth/helpers/apperrors"
	"backend/modules/auth/models"
	"log"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// AuthAccount verifies that the Authorization header is a valid JWT from the server
func AuthAccount(service models.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := authHeader{}

		// bind the incoming authorization header and check for errors
		if err := ctx.ShouldBindHeader(&header); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				err := apperrors.NewBadRequest("Invalid request paramaters. See invalidArgs")
				ctx.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				ctx.Abort()
				return
			}

			//
			err := apperrors.NewInternal()
			ctx.JSON(err.Status(), gin.H{
				"error": err,
			})
			ctx.Abort()
			return
		}

		authHeader := strings.Split(header.IDToken, "Bearer ")
		log.Print("Hello Auth:Bearer token from middleware success") //authHeader
		if len(authHeader) < 2 {
			err := apperrors.NewAuthorization("Must provide an Authorization header with format `Bearer: {token}`")
			ctx.JSON(err.Status(), gin.H{
				"error": err,
			})

			ctx.Abort()
			return
		}

		// Perform IDToken validation here
		idToken, err := service.ValidateIDToken(authHeader[1])
		if err != nil {
			err := apperrors.NewAuthorization("debug middleware: Provided Authorization token is invalid")
			ctx.JSON(err.Status(), gin.H{
				"error": err,
			})
			ctx.Abort()
			return
		}

		ctx.Set("account", idToken)
		ctx.Next()
	}
}
