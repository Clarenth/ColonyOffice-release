package handlers

import (
	"backend/internal/helpers/apperrors"
	"backend/internal/mocks"
	"backend/internal/models"
	"bytes"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignup(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Email and Password Required", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockEmployeeService := new(mocks.MockEmployeeService)
		mockEmployeeService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*models.EmployeeAccountData")).Return(nil)

		// a response recorder for getting written http response
		responseRecorder := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		Routes(&HandlerConfig{
			Router:          router,
			EmployeeService: mockEmployeeService,
		})

		// create a request body with empty email and password
		requestBody, err := json.Marshal(gin.H{
			"email":    "",
			"password": "",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, 400, responseRecorder.Code)
		mockEmployeeService.AssertNotCalled(t, "Signup")
	})

	t.Run("Invalid email", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockEmployeeService := new(mocks.MockEmployeeService)
		mockEmployeeService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*models.EmployeeAccountData")).Return(nil)

		// a response recorder for getting written http response
		responseRecorder := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		Routes(&HandlerConfig{
			Router:          router,
			EmployeeService: mockEmployeeService,
		})

		// create a request body with empty email and password
		requestBody, err := json.Marshal(gin.H{
			"email":    "crag@tarr",
			"password": "1234567890-==-0987654321",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, 400, responseRecorder.Code)
		mockEmployeeService.AssertNotCalled(t, "Signup")
	})

	t.Run("Password too short", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockEmployeeService := new(mocks.MockEmployeeService)
		mockEmployeeService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*models.EmployeeAccountData")).Return(nil)

		// a response recorder for getting written http response
		responseRecorder := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		Routes(&HandlerConfig{
			Router:          router,
			EmployeeService: mockEmployeeService,
		})

		// create a request body with empty email and password
		requestBody, err := json.Marshal(gin.H{
			"email":    "crag@tarr.tr",
			"password": "1234",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, 400, responseRecorder.Code)
		mockEmployeeService.AssertNotCalled(t, "Signup")
	})
	t.Run("Password too long", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockEmployeeService := new(mocks.MockEmployeeService)
		mockEmployeeService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*models.EmployeeAccountData")).Return(nil)

		// a response recorder for getting written http response
		responseRecorder := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		Routes(&HandlerConfig{
			Router:          router,
			EmployeeService: mockEmployeeService,
		})

		// create a request body with empty email and password
		requestBody, err := json.Marshal(gin.H{
			"email":    "bob@bob.com",
			"password": "1234567890-==-0987654321dh82385h59jmfnve0p5kgndofe0tk540-heiohj5jh-y5096jhu560-yi56uj4r5jy-huj509ohkup5rjhmu45mhuk4,ju,4r50-4=-6uk,jnr[p83rhfejiew9ihui38u4e94tnfer93htopgngl;;r03-4ngvge84jntg9ektntg0e3pmt",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, 400, responseRecorder.Code)
		mockEmployeeService.AssertNotCalled(t, "Signup")
	})

	t.Run("Error returned from EmployeeService", func(t *testing.T) {
		user := &models.EmployeeAccountData{
			Email:    "crag@tarr.tr",
			Password: "1234567890-==-0987654321",
		}

		mockEmployeeService := new(mocks.MockEmployeeService)
		mockEmployeeService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), user).Return(apperrors.NewConflict("User Already Exists", user.Email))

		// a response recorder for getting written http response
		responseRecorder := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		Routes(&HandlerConfig{
			Router:          router,
			EmployeeService: mockEmployeeService,
		})

		// create a request body with empty email and password
		requestBody, err := json.Marshal(gin.H{
			"email":    user.Email,
			"password": user.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, 409, responseRecorder.Code)
		mockEmployeeService.AssertExpectations(t)
	})

	t.Run("Successful Token Creation", func(t *testing.T) {
		user := &models.EmployeeAccountData{
			Email:    "crag@tarr.tr",
			Password: "1234567890-==-0987654321",
		}

		mockTokenResponse := &models.TokenPair{
			IDToken:      "idToken",
			RefreshToken: "refreshToken",
		}

		mockEmployeeService := new(mocks.MockEmployeeService)
		mockTokenService := new(mocks.MockTokenService)

		mockEmployeeService.
			On("Signup", mock.AnythingOfType("*context.emptyCtx"), user).
			Return(nil)
		mockTokenService.
			On("NewTokenPairFromUser", mock.AnythingOfType("*context.emptyCtx"), user, "").
			Return(mockTokenResponse, nil)

		// a response recorder for getting written http response
		responseRecorder := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		Routes(&HandlerConfig{
			Router:          router,
			EmployeeService: mockEmployeeService,
			TokenService:    mockTokenService,
		})

		// create a request body with empty email and password
		requestBody, err := json.Marshal(gin.H{
			"email":    user.Email,
			"password": user.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		respBody, err := json.Marshal(gin.H{
			"tokens": mockTokenResponse,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, responseRecorder.Code)
		assert.Equal(t, respBody, responseRecorder.Body.Bytes())

		mockEmployeeService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("Failed Token Creation", func(t *testing.T) {
		user := &models.EmployeeAccountData{
			Email:    "crag@tar.tr",
			Password: "1234567890-==-0987654321",
		}

		mockErrorResponse := apperrors.NewInternal()

		mockEmployeeService := new(mocks.MockEmployeeService)
		mockTokenService := new(mocks.MockTokenService)

		mockEmployeeService.
			On("Signup", mock.AnythingOfType("*context.emptyCtx"), user).
			Return(nil)
		mockTokenService.
			On("NewTokenPairFromUser", mock.AnythingOfType("*context.emptyCtx"), user, "").
			Return(nil, mockErrorResponse)

		// a response recorder for getting written http response
		responseRecorder := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		Routes(&HandlerConfig{
			Router:          router,
			EmployeeService: mockEmployeeService,
			TokenService:    mockTokenService,
		})

		// create a request body with empty email and password
		requestBody, err := json.Marshal(gin.H{
			"email":    user.Email,
			"password": user.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockErrorResponse,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockErrorResponse.Status(), responseRecorder.Code)
		assert.Equal(t, respBody, responseRecorder.Body.Bytes())

		mockEmployeeService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})
}
