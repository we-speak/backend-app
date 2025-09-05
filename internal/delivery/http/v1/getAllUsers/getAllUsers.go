package getAllUsers

import (
	"backend-app/internal/storage/models"
	"backend-app/pkg/api/response"
	"backend-app/pkg/sl"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Getter interface {
	GetAllUsers(offset int, limit int) ([]models.User, error)
}

// New godoc
// @Summary Get all users
// @Description Returns a list of all users with pagination
// @Tags users
// @Produce json
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Success 200 {array} models.User
// @Failure 422 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/user/all [get]
func New(log *slog.Logger, getter Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.GetAllUsers"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			log.Error("invalid query parameter", sl.Error(err))
			render.Status(r, http.StatusUnprocessableEntity)
			return
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			log.Error("invalid query parameter", sl.Error(err))
			render.Status(r, http.StatusUnprocessableEntity)
			return
		}
		users, err := getter.GetAllUsers(offsetInt, limitInt)
		if err != nil {
			log.Error("failed to get users", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get users"))
			return
		}

		log.Info("users retrieved successfully", "count", len(users))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, users)
	}
}
