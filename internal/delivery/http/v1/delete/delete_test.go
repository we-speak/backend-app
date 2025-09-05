package delete_test

import (
	delete2 "backend-app/internal/delivery/http/v1/delete"
	"backend-app/pkg/api/response"
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

type mockDeleter struct {
	DeleteFn func(id uint) error
}

func (m *mockDeleter) DeleteUser(id uint) error {
	return m.DeleteFn(id)
}

func TestDeleteUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		urlParam       string
		mockDeleteErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "invalid_id",
			urlParam:       "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid user id",
		},
		{
			name:           "user_not_found",
			urlParam:       "123",
			mockDeleteErr:  gorm.ErrRecordNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "user not found",
		},
		{
			name:           "internal_error",
			urlParam:       "123",
			mockDeleteErr:  errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to delete user",
		},
		{
			name:           "success",
			urlParam:       "123",
			expectedStatus: http.StatusOK,
			expectedBody:   "", // will check status OK
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.urlParam, nil)
			rr := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Use(middleware.RequestID)
			router.Use(render.SetContentType(render.ContentTypeJSON))
			router.Delete("/users/{id}", delete2.New(slog.Default(), &mockDeleter{
				DeleteFn: func(id uint) error {
					return tt.mockDeleteErr
				},
			}))

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var res response.Response
			_ = render.DecodeJSON(rr.Body, &res)

			if tt.expectedBody != "" {
				assert.Equal(t, "Error", res.Status)
				assert.Equal(t, tt.expectedBody, res.Error)
			} else {
				assert.Equal(t, "OK", res.Status)
			}
		})
	}
}
