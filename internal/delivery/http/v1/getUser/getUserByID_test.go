package getUser_test

import (
	v1Router "backend-app/internal/delivery/http/v1/getUser"
	"backend-app/internal/storage/models"
	_ "bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// mockStorage implements the method needed for testing
type mockStorage struct {
	user *models.User
	err  error
}

func (m *mockStorage) GetUserByID(id uint) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name           string
		paramID        string
		mockUser       *models.User
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "invalid id",
			paramID:        "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"Error","error":"invalid user id"}`,
		},
		{
			name:           "user not found",
			paramID:        "1",
			mockError:      gorm.ErrRecordNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":"Error","error":"user not found"}`,
		},
		{
			name:           "internal error",
			paramID:        "1",
			mockError:      errors.New("db failure"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"status":"Error","error":"failed to get user"}`,
		},
		{
			name:           "success",
			paramID:        "1",
			mockUser:       &models.User{ID: 1, Email: "test@example.com"},
			expectedStatus: http.StatusOK,
			expectedBody: `{
		"id": 1,
		"username": "",
		"password": "",
		"email": "test@example.com",
		"role": "",
		"country": "",
		"createdAt": "0001-01-01T00:00:00Z",
		"updatedAt": "0001-01-01T00:00:00Z"
	}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &mockStorage{
				user: tt.mockUser,
				err:  tt.mockError,
			}
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.paramID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler := v1Router.New(slog.Default(), storage)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}
		})

	}
}
