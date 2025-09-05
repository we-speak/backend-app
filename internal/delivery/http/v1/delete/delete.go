package delete

import (
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

type deleter interface {
	DeleteUser(id uint) error
}

// New godoc
// @Summary Delete user
// @Description Deletes a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/user/{id} [delete]
func New(log *slog.Logger, deleter deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.DeleteUser"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		idUint, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			log.Error("invalid user id", "param", idStr, "error", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid user id"))
			return
		}

		err = deleter.DeleteUser(uint(idUint))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("user not found", "id", idUint)
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("user not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete user", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to delete user"))
			return
		}

		log.Info("user deleted successfully", "id", idUint)
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.OK())
	}
}
