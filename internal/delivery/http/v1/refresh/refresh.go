package refresh

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
	"github.com/golang-jwt/jwt/v5"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// New godoc
// @Summary Refresh token pair
// @Description Generates new access and refresh tokens using valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body refresh.RefreshRequest true "Refresh token"
// @Success 200 {object} config.TokenPair
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/refresh [post]
func New(log *slog.Logger, storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var request RefreshRequest
		if err := render.DecodeJSON(r.Body, &request); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid request"})
			return
		}

		// Парсим refresh token
		token, err := jwt.ParseWithClaims(request.RefreshToken, &config.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return config.RefreshJWTSecret, nil
		})

		if err != nil || !token.Valid {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Invalid refresh token"})
			return
		}

		claims, ok := token.Claims.(*config.Claims)
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Invalid token claims"})
			return
		}

		tokenPair, err := generator.GenerateTokenPair(claims.UserID, claims.Role)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Could not generate tokens"})
			return
		}
		user, err := storage.GetUserByID(claims.UserID)
		if err != nil {
			log.Error("err", sl.Error(err))
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		if request.RefreshToken != user.RefreshToken {
			log.Error("err", "refresh token закончился хули")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid refresh token"))
			return
		}

		user.RefreshToken = tokenPair.RefreshToken
		user.TokenExpiry = time.Now().Add(config.RefreshTokenExpiry)

		if err := storage.UpdateUser(user); err != nil {
			log.Error("err", sl.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.JSON(w, r, tokenPair)
	}
}
