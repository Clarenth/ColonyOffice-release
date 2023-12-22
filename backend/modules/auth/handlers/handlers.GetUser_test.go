package handlers

import (
	"backend/internal/helpers/apperrors"
	"backend/internal/mocks"
	"backend/internal/models"

	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCurrentUser(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockEmployeeResponse := &models.EmployeeAccountData{
			IDCode: uid,
			Email:  "crag@tarr.gov",
			EmployeeIdentityData: &models.EmployeeIdentityData{
				FirstName: "Crag Tarr",
			},
		}

		mockEmployeeService := new(mocks.MockEmployeeService)
		mockEmployeeService.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), uid).Return(mockEmployeeResponse, nil)

		// a response recorder for getting http response
		responseRecorder := httptest.NewRecorder()

		// use a middleware to set context for test
		// the only claims we care about in this test is the UID
		router := gin.Default()
		router.Use(func(ctx *gin.Context) {
			ctx.Set("account", &models.EmployeeAccountData{
				IDCode: uid,
			})
		})

		Routes(&HandlerConfig{
			Router:          router,
			EmployeeService: mockEmployeeService,
		})
		request, err := http.NewRequest(http.MethodGet, "/auth/user", nil)
		assert.NoError(t, err)

		router.ServeHTTP(responseRecorder, request)

		responseBody, err := json.Marshal(gin.H{
			"account": mockEmployeeResponse,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Equal(t, responseBody, responseRecorder.Body.Bytes())
		mockEmployeeService.AssertExpectations(t) // asserts that EmployeeService was called
	})

	t.Run("NoContextUser", func(t *testing.T) {
		mockEmployeeService := new(mocks.MockEmployeeService)
		mockEmployeeService.On("GetUser", mock.Anything, mock.Anything).Return(nil, nil)

		rr := httptest.NewRecorder()

		router := gin.Default()
		Routes(&HandlerConfig{
			Router:          router,
			EmployeeService: mockEmployeeService,
		})

		request, err := http.NewRequest(http.MethodGet, "/auth/user", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockEmployeeService.AssertNotCalled(t, "GetUser", mock.Anything)
	})

	t.Run("NotFound", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockEmployeeService := new(mocks.MockEmployeeService)
		mockEmployeeService.On("GetUser", mock.Anything, uid).Return(nil, fmt.Errorf("An error occured down the call chain"))

		// Response Recorder ofr http responses
		responseRecorder := httptest.NewRecorder()

		router := gin.Default()
		router.Use(func(ctx *gin.Context) {
			ctx.Set("account", &models.EmployeeAccountData{
				IDCode: uid,
			},
			)
		})

		Routes(&HandlerConfig{
			Router:          router,
			EmployeeService: mockEmployeeService,
		})

		request, err := http.NewRequest(http.MethodGet, "/auth/user", nil)
		assert.NoError(t, err)

		router.ServeHTTP(responseRecorder, request)

		responseError := apperrors.NewNotFound("account", uid.String())

		responseBody, err := json.Marshal(gin.H{
			"error": responseError,
		})
		assert.NoError(t, err)

		assert.Equal(t, responseError.Status(), responseRecorder.Code)
		assert.Equal(t, responseBody, responseRecorder.Body.Bytes())
		mockEmployeeService.AssertExpectations(t)
	})
}
