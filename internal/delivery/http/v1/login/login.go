package login

import (
	"backend-app/internal/config"
	"backend-app/internal/storage/postgres"
	"backend-app/pkg/api/response"
	"backend-app/pkg/jwt/generator"
	"backend-app/pkg/sl"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// New godoc
// @Summary Login
// @Description Authenticates user and returns token pair
// @Tags auth
// @Accept json
// @Produce json
// @Param input body login.LoginRequest true "Credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/login [post]
func New(log *slog.Logger, storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials LoginRequest

		if err := render.DecodeJSON(r.Body, &credentials); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid request"})
			return
		}

		user, err := storage.GetUserByUsername(credentials.Username)
		if err != nil {
			log.Error("invalid request", sl.Error(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid username"})
		}

		if err := user.CheckPassword(credentials.Password); err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Invalid credentials"})
			return
		}
		tokens, err := generator.GenerateTokenPair(user.ID, user.Role)
		if err != nil {
			log.Error("error", sl.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create token pair"))
		}

		user.RefreshToken = tokens.RefreshToken
		user.TokenExpiry = time.Now().Add(config.RefreshTokenExpiry)

		if err := storage.UpdateUser(user); err != nil {
			log.Error("err", sl.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, map[string]string{
			"acess_token":   tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		})
	}
}
