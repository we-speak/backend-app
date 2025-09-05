package server

import (
	"backend-app/internal/config"
	"log/slog"
	"net/http"
)

func ListenAndServe(router http.Handler, cfg *config.Config) error {
	server := &http.Server{
		Addr:         cfg.HTTPServer.Host,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	slog.Info("Starting server", "on", cfg.HTTPServer.Host)
	return server.ListenAndServe()
}
