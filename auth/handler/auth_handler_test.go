package handler

import (
	"7solutions-challenge/auth/domain"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(
	ctx context.Context,
	req domain.UserRequest,
) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockAuthService) Login(
	ctx context.Context,
	req domain.LoginRequest,
) (string, error) {
	args := m.Called(ctx, req)

	return args.String(0), args.Error(1)
}

func setupRouter(service *mockAuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	h := NewAuthHandler(service)

	router.POST("/register", h.Register)
	router.POST("/login", h.Login)

	return router
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		setupMock      func(*mockAuthService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			body: `{
				"name": "John Doe",
				"email": "john@example.com",
				"password": "password123"
			}`,
			setupMock: func(m *mockAuthService) {
				m.On(
					"Register",
					mock.Anything,
					mock.AnythingOfType("domain.UserRequest"),
				).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "user registered successfully",
		},
		{
			name: "invalid json",
			body: `{
				"name": "John Doe",
				"email":
			`,
			setupMock:      func(m *mockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "error",
		},
		{
			name: "validation error - missing name",
			body: `{
				"name": "",
				"email": "john@example.com",
				"password": "password123"
			}`,
			setupMock:      func(m *mockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "UserRequest.Name",
		},
		{
			name: "validation error - invalid email",
			body: `{
				"name": "John Doe",
				"email": "invalid-email",
				"password": "password123"
			}`,
			setupMock:      func(m *mockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "UserRequest.Email",
		},
		{
			name: "validation error - missing password",
			body: `{
				"name": "John Doe",
				"email": "john@example.com",
				"password": ""
			}`,
			setupMock:      func(m *mockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "UserRequest.Password",
		},
		{
			name: "validation error - name too long",
			body: `{
				"name": "` + strings.Repeat("a", 301) + `",
				"email": "john@example.com",
				"password": "password123"
			}`,
			setupMock:      func(m *mockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "UserRequest.Name",
		},
		{
			name: "validation error - password too long",
			body: `{
				"name": "John Doe",
				"email": "john@example.com",
				"password": "` + strings.Repeat("a", 73) + `"
			}`,
			setupMock:      func(m *mockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "UserRequest.Password",
		},
		{
			name: "service error",
			body: `{
				"name": "John Doe",
				"email": "john@example.com",
				"password": "password123"
			}`,
			setupMock: func(m *mockAuthService) {
				m.On(
					"Register",
					mock.Anything,
					mock.AnythingOfType("domain.UserRequest"),
				).Return(errors.New("registration failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "registration failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(mockAuthService)
			tt.setupMock(service)

			router := setupRouter(service)

			req := httptest.NewRequest(
				http.MethodPost,
				"/register",
				strings.NewReader(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectedBody)

			service.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		setupMock      func(*mockAuthService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			body: `{
				"email": "john@example.com",
				"password": "password123"
			}`,
			setupMock: func(m *mockAuthService) {
				m.On(
					"Login",
					mock.Anything,
					mock.AnythingOfType("domain.LoginRequest"),
				).Return("test-token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "test-token",
		},
		{
			name: "invalid json",
			body: `{
				"email":
			`,
			setupMock:      func(m *mockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "error",
		},
		{
			name: "validation error - missing email",
			body: `{
				"email": "",
				"password": "password123"
			}`,
			setupMock:      func(m *mockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "LoginRequest.Email",
		},
		{
			name: "validation error - invalid email",
			body: `{
				"email": "invalid-email",
				"password": "password123"
			}`,
			setupMock:      func(m *mockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "LoginRequest.Email",
		},
		{
			name: "validation error - missing password",
			body: `{
				"email": "john@example.com",
				"password": ""
			}`,
			setupMock:      func(m *mockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "LoginRequest.Password",
		},
		{
			name: "validation error - password too long",
			body: `{
				"email": "john@example.com",
				"password": "` + strings.Repeat("a", 73) + `"
			}`,
			setupMock:      func(m *mockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "LoginRequest.Password",
		},
		{
			name: "invalid credentials",
			body: `{
				"email": "john@example.com",
				"password": "wrongpassword"
			}`,
			setupMock: func(m *mockAuthService) {
				m.On(
					"Login",
					mock.Anything,
					mock.AnythingOfType("domain.LoginRequest"),
				).Return("", errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(mockAuthService)
			tt.setupMock(service)

			router := setupRouter(service)

			req := httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectedBody)

			service.AssertExpectations(t)
		})
	}
}
