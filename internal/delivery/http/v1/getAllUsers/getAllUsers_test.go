package getAllUsers_test

import (
	"backend-app/internal/delivery/http/v1/getAllUsers"
	"backend-app/internal/storage/models"
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
)

type mockGetter struct {
	GetAllUsersFn func() ([]models.User, error)
}

func (m *mockGetter) GetAllUsers() ([]models.User, error) {
	return m.GetAllUsersFn()
}

func TestGetAllUsersHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockReturn     []models.User
		mockError      error
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "internal_error",
			mockError:      errors.New("db failure"),
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name: "success",
			mockReturn: []models.User{
				{ID: 1, Email: "test1@example.com"},
				{ID: 2, Email: "test2@example.com"},
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			rr := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Use(middleware.RequestID)
			router.Use(render.SetContentType(render.ContentTypeJSON))
			router.Get("/users", getAllUsers.New(slog.Default(), &mockGetter{
				GetAllUsersFn: func() ([]models.User, error) { //FIXME: попозже пофикшу здесь надо добавить оффсет и лимит
					return tt.mockReturn, tt.mockError
				},
			}))

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectError {
				var res struct {
					Status string `json:"status"`
					Error  string `json:"error"`
				}
				err := json.Unmarshal(rr.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Equal(t, "Error", res.Status)
				assert.NotEmpty(t, res.Error)
			} else {
				var users []models.User
				err := json.Unmarshal(rr.Body.Bytes(), &users)
				assert.NoError(t, err)
				assert.Len(t, users, len(tt.mockReturn))
			}
		})
	}
}
