package authMiddleware

import (
	"backend-app/internal/config"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

var (
	AccessTokenAuth  *jwtauth.JWTAuth
	RefreshTokenAuth *jwtauth.JWTAuth
)

func init() {
	AccessTokenAuth = jwtauth.New("HS256", config.JWTSecret, nil)
	RefreshTokenAuth = jwtauth.New("HS256", config.RefreshJWTSecret, nil)
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil || token == nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Not found"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Not found"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
