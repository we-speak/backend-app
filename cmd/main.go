// @title DishFinder auth service docs
// @version 1.0
// @description
// @host localhost:8080
// @BasePath /
package main

import (
	"backend-app/internal/config"
	router "backend-app/internal/delivery/http"
	"backend-app/internal/server"
	"backend-app/internal/storage/postgres"
	"backend-app/pkg/logger"
	"backend-app/pkg/sl"
	"log"
	"log/slog"
)

func main() {

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	log := logger.New(cfg.Env)
	storage, err := postgres.New(cfg)

	if err != nil {
		log.Error("Error connect to postgreSQL", sl.Error(err))
	}

	log.Info("Starting server", "env", cfg.Env, "host", cfg.HTTPServer.Host)
	log.Info("Server timeout", "timeout", cfg.HTTPServer.Timeout)
	log.Info("Server idle timeout", "idle_timeout", cfg.HTTPServer.IdleTimeout)
	r := router.InitRoutes(log, &storage, cfg)

	if err := server.ListenAndServe(r, cfg); err != nil {
		log.Error("Error starting server: %v", slog.String("err", err.Error()))
	}
}
