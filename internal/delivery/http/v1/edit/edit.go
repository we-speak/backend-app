package edit

import (
	"backend-app/internal/storage/models"
	"backend-app/pkg/api/response"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Updater interface {
	UpdateUser(user *models.User) error
}

// New godoc
// @Summary Update user
// @Description Updates user data
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.User true "User data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/user [put]
func New(log *slog.Logger, updater Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.UpdateUser"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.User
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", "error", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}
		log.Info("request body decoded", slog.Any("user", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("validation failed", "error", err)
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, response.Error("validation failed"))
			return
		}

		err := updater.UpdateUser(&req)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("user not found", "id", req.ID)
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("user not found"))
			return
		}
		if err != nil {
			log.Error("failed to update user", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to update user"))
			return
		}

		log.Info("user updated successfully", "id", req.ID)
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.OK())
	}
}
