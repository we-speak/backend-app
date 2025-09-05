package edit_test

import (
	"backend-app/internal/delivery/http/v1/edit"
	"backend-app/internal/storage/models"
	"backend-app/pkg/api/response"
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockUpdater struct {
	UpdateFn func(user *models.User) error
}

func (m *mockUpdater) UpdateUser(user *models.User) error {
	return m.UpdateFn(user)
}

func TestUpdateUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockUpdateErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body",
		},
		{
			name: "validation error",
			requestBody: map[string]interface{}{
				"id":    1,
				"email": "", // required field is empty
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "validation failed",
		},
		{
			name: "user not found",
			requestBody: models.User{
				ID:       1,
				Username: "updated",
				Password: "newpass",
				Email:    "user@example.com",
				Role:     "user",
				Country:  "RU",
			},
			mockUpdateErr:  gorm.ErrRecordNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "user not found",
		},
		{
			name: "internal error",
			requestBody: models.User{
				ID:       1,
				Username: "updated",
				Password: "newpass",
				Email:    "user@example.com",
				Role:     "user",
				Country:  "RU",
			},
			mockUpdateErr:  errors.New("db failure"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to update user",
		},
		{
			name: "success",
			requestBody: models.User{
				ID:       1,
				Username: "updated",
				Password: "newpass",
				Email:    "user@example.com",
				Role:     "user",
				Country:  "RU",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "", // should return OK response
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			switch v := tt.requestBody.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, err = json.Marshal(v)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler := edit.New(slog.Default(), &mockUpdater{
				UpdateFn: func(user *models.User) error {
					return tt.mockUpdateErr
				},
			})

			r := chi.NewRouter()
			r.Use(middleware.RequestID)
			r.Use(render.SetContentType(render.ContentTypeJSON))
			r.Put("/users", handler)

			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				var res response.Response
				err := json.Unmarshal(rr.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Equal(t, "Error", res.Status)
				assert.Equal(t, tt.expectedBody, res.Error)
			} else {
				var res response.Response
				err := json.Unmarshal(rr.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Equal(t, "OK", res.Status)
			}
		})
	}
}
