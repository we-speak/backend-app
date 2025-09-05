package v1Router

import (
	authMiddleware "backend-app/internal/delivery/http/middleware/auth"
	delete2 "backend-app/internal/delivery/http/v1/delete"
	"backend-app/internal/delivery/http/v1/edit"
	"backend-app/internal/delivery/http/v1/getAllUsers"
	"backend-app/internal/delivery/http/v1/getUser"
	"backend-app/internal/delivery/http/v1/login"
	"backend-app/internal/delivery/http/v1/refresh"
	"backend-app/internal/delivery/http/v1/register"
	"backend-app/internal/storage/postgres"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func New(log *slog.Logger, storage *postgres.Storage) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Throttle(100))

	r.Post("/refresh", refresh.New(log, storage))
	r.Post("/register", register.New(log, storage))
	r.Group(func(r chi.Router) {

		r.Use(jwtauth.Verifier(authMiddleware.RefreshTokenAuth))

		r.Use(authMiddleware.Authenticator)

	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(authMiddleware.AccessTokenAuth))
		r.Use(authMiddleware.Authenticator)
		r.Use(authMiddleware.AdminOnly)

		r.Delete("/user/{id}", delete2.New(log, storage))
		r.Get("/user/all", getAllUsers.New(log, storage))
		r.Get("/user/{id}", getUser.New(log, storage))
		r.Put("/user", edit.New(log, storage))

	})
	r.Post("/login", login.New(log, storage))
	return r
}
