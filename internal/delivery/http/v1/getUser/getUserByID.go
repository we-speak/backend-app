package getUser

import (
	"backend-app/internal/storage/models"
	"backend-app/pkg/api/response"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

type storage interface {
	GetUserByID(id uint) (*models.User, error)
}

// New godoc
// @Summary Get user by ID
// @Description Returns user data by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/user/{id} [get]
func New(log *slog.Logger, storage storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.GetUserByID"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			log.Error("invalid user id", "param", idParam, "error", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid user id"))
			return
		}

		user, err := storage.GetUserByID(uint(id))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("user not found", "id", id)
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("user not found"))
			return
		}
		if err != nil {
			log.Error("failed to get user", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get user"))
			return
		}

		log.Info("user retrieved successfully", slog.Any("user", user))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, user)
	}
}
