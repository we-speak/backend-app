package register

import (
	"backend-app/internal/storage/models"
	"backend-app/pkg/api/response"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	validator2 "github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Saver interface {
	CreateUser(user *models.User) error
}

// New godoc
// @Summary Register new user
// @Description Create user account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.User true "User data"
// @Success 201 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /v1/register [post]
func New(log *slog.Logger, saver Saver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.New"
		response.OK()
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.User
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode body", "error", err)
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, response.Error("failed to decode body"))
			return
		}
		log.Info("request body decoded", slog.Any("body", req))

		if err := validator2.New().Struct(req); err != nil {
			log.Error("failed to validate body", "error", err)
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, response.Error("failed to validate body"))
			return
		}
		if err := req.HashPassword(); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Could not hash password"})
			return
		}

		err = saver.CreateUser(&req)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Info("user already exists", "error", err)
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, response.Error("user already exists"))

			return
		}
		if err != nil {
			log.Error("failed to create user", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create user"))
			return
		}

		log.Info("user created successfully")
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response.OK())
	}
}
