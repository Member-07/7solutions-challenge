package handler

import (
	"7solutions-challenge/user/domain"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// Mock User Service
// ============================================================

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Create(
	ctx context.Context,
	req domain.UserRequest,
) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockUserService) GetByID(
	ctx context.Context,
	id string,
) (domain.UserResponse, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(domain.UserResponse), args.Error(1)
}

func (m *mockUserService) List(
	ctx context.Context,
) ([]domain.UserResponse, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]domain.UserResponse), args.Error(1)
}

func (m *mockUserService) Update(
	ctx context.Context,
	req domain.UserUpdateRequest,
) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockUserService) Delete(
	ctx context.Context,
	id string,
) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUserService) Count(
	ctx context.Context,
) (int64, error) {
	args := m.Called(ctx)

	return args.Get(0).(int64), args.Error(1)
}

// ============================================================
// Router
// ============================================================

func setupRouter(service *mockUserService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	h := NewUserHandler(service)

	router.POST("/users", h.Create)
	router.GET("/users/count", h.Count)
	router.GET("/users/:id", h.GetByID)
	router.GET("/users", h.List)
	router.PUT("/users/:id", h.Update)
	router.DELETE("/users/:id", h.Delete)

	return router
}

// ============================================================
// Create
// ============================================================

func TestUserHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		setupMock      func(*mockUserService)
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
			setupMock: func(m *mockUserService) {
				m.On(
					"Create",
					mock.Anything,
					mock.AnythingOfType("domain.UserRequest"),
				).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "user created successfully",
		},
		{
			name: "invalid json",
			body: `{
				"name": "John Doe",
				"email":
			`,
			setupMock:      func(m *mockUserService) {},
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
			setupMock:      func(m *mockUserService) {},
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
			setupMock:      func(m *mockUserService) {},
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
			setupMock:      func(m *mockUserService) {},
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
			setupMock:      func(m *mockUserService) {},
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
			setupMock:      func(m *mockUserService) {},
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
			setupMock: func(m *mockUserService) {
				m.On(
					"Create",
					mock.Anything,
					mock.AnythingOfType("domain.UserRequest"),
				).Return(errors.New("create user failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "create user failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(mockUserService)
			tt.setupMock(service)

			router := setupRouter(service)

			req := httptest.NewRequest(
				http.MethodPost,
				"/users",
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

// ============================================================
// GetByID
// ============================================================

func TestUserHandler_GetByID(t *testing.T) {
	user := domain.UserResponse{
		ID:        "user-123",
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		service := new(mockUserService)

		service.On(
			"GetByID",
			mock.Anything,
			"user-123",
		).Return(user, nil)

		router := setupRouter(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/user-123",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"id":"user-123"`)
		assert.Contains(t, rec.Body.String(), `"name":"John Doe"`)
		assert.Contains(t, rec.Body.String(), `"email":"john@example.com"`)

		service.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		service := new(mockUserService)

		service.On(
			"GetByID",
			mock.Anything,
			"user-123",
		).Return(
			domain.UserResponse{},
			errors.New("user not found"),
		)

		router := setupRouter(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/user-123",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "user not found")

		service.AssertExpectations(t)
	})
}

// ============================================================
// List
// ============================================================

func TestUserHandler_List(t *testing.T) {
	users := []domain.UserResponse{
		{
			ID:        "user-1",
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
		},
		{
			ID:        "user-2",
			Name:      "Jane Doe",
			Email:     "jane@example.com",
			CreatedAt: time.Now(),
		},
	}

	t.Run("success", func(t *testing.T) {
		service := new(mockUserService)

		service.On(
			"List",
			mock.Anything,
		).Return(users, nil)

		router := setupRouter(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"user-1"`)
		assert.Contains(t, rec.Body.String(), `"user-2"`)

		service.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		service := new(mockUserService)

		service.On(
			"List",
			mock.Anything,
		).Return(
			nil,
			errors.New("failed to get users"),
		)

		router := setupRouter(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to get users")

		service.AssertExpectations(t)
	})
}

// ============================================================
// Update
// ============================================================

func TestUserHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		body           string
		setupMock      func(*mockUserService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			id:   "user-123",
			body: `{
				"name": "John Updated",
				"email": "john.updated@example.com"
			}`,
			setupMock: func(m *mockUserService) {
				m.On(
					"Update",
					mock.Anything,
					mock.AnythingOfType("domain.UserUpdateRequest"),
				).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "user updated successfully",
		},
		{
			name: "invalid json",
			id:   "user-123",
			body: `{
				"name":
			`,
			setupMock:      func(m *mockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "error",
		},
		{
			name: "validation error - missing name",
			id:   "user-123",
			body: `{
				"name": "",
				"email": "john@example.com"
			}`,
			setupMock:      func(m *mockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "UserUpdateRequest.Name",
		},
		{
			name: "validation error - invalid email",
			id:   "user-123",
			body: `{
				"name": "John Doe",
				"email": "invalid-email"
			}`,
			setupMock:      func(m *mockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "UserUpdateRequest.Email",
		},
		{
			name: "validation error - name too long",
			id:   "user-123",
			body: `{
				"name": "` + strings.Repeat("a", 301) + `",
				"email": "john@example.com"
			}`,
			setupMock:      func(m *mockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "UserUpdateRequest.Name",
		},
		{
			name: "service error",
			id:   "user-123",
			body: `{
				"name": "John Updated",
				"email": "john.updated@example.com"
			}`,
			setupMock: func(m *mockUserService) {
				m.On(
					"Update",
					mock.Anything,
					mock.AnythingOfType("domain.UserUpdateRequest"),
				).Return(errors.New("update user failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "update user failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(mockUserService)
			tt.setupMock(service)

			router := setupRouter(service)

			req := httptest.NewRequest(
				http.MethodPut,
				"/users/"+tt.id,
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

// ============================================================
// Delete
// ============================================================

func TestUserHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(mockUserService)

		service.On(
			"Delete",
			mock.Anything,
			"user-123",
		).Return(nil)

		router := setupRouter(service)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/users/user-123",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())

		service.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		service := new(mockUserService)

		service.On(
			"Delete",
			mock.Anything,
			"user-123",
		).Return(errors.New("delete user failed"))

		router := setupRouter(service)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/users/user-123",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "delete user failed")

		service.AssertExpectations(t)
	})
}

// ============================================================
// Count
// ============================================================

func TestUserHandler_Count(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(mockUserService)

		service.On(
			"Count",
			mock.Anything,
		).Return(int64(10), nil)

		router := setupRouter(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/count",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"count":10`)

		service.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		service := new(mockUserService)

		service.On(
			"Count",
			mock.Anything,
		).Return(
			int64(0),
			errors.New("count users failed"),
		)

		router := setupRouter(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/count",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "count users failed")

		service.AssertExpectations(t)
	})
}
